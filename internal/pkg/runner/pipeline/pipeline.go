/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package pipeline

import (
	"context"
	"fmt"
	"go-brunel/internal/pkg/runner/recorder"
	"go-brunel/internal/pkg/runner/runtime"
	"go-brunel/internal/pkg/shared"
	"go-brunel/internal/pkg/shared/util"
	"log"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

var (
	defaultWaitTimeout = 30 * time.Second
)

type Pipeline struct {
	Runtime  runtime.Runtime
	Recorder recorder.Recorder
}

func (pipeline *Pipeline) cleanUp(context context.Context, containerIDs []shared.ContainerID) error {
	var err error
	for _, containerID := range containerIDs {
		if e := pipeline.Runtime.TerminateContainer(context, containerID); e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, "error terminating container"))
		}

		if e := pipeline.Recorder.RecordContainerState(containerID, shared.ContainerStateStopped); e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, "error updating terminating container status"))
		}
	}
	return err
}

func (pipeline *Pipeline) executeStage(context context.Context, jobID shared.JobID, stageID shared.StageID, stage shared.Stage) ([]shared.ContainerID, error) {

	// We use this to return any containers that need cleaned up during an error
	var containerIDs []shared.ContainerID

	// If we have services, dispatch them
	if stage.Services != nil {
		for _, sidecar := range stage.Services {
			// Dispatch the container, it may not be started/stopWaiting at this point
			containerID, err := pipeline.Runtime.DispatchContainer(context, jobID, sidecar)
			if containerID != shared.EmptyContainerID {
				containerIDs = append(containerIDs, containerID)
			}
			if err != nil {
				return containerIDs, errors.Wrap(err, "error dispatching sidecar service container")
			}

			// Record our container as starting
			err = pipeline.Recorder.RecordContainer(jobID, containerID, shared.ContainerMeta{StageID: stageID, Service: true}, sidecar, shared.ContainerStateStarting)
			if err != nil {
				return containerIDs, errors.Wrap(err, "error recording sidecar service container creation")
			}

			// Wait for our container to be running, then mark it as running
			if err = pipeline.Runtime.WaitForContainer(context, containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitRunning}); err != nil {
				return containerIDs, errors.Wrap(err, "error waiting for sidecar service container to be running")
			}

			if err = pipeline.Recorder.RecordContainerState(containerID, shared.ContainerStateRunning); err != nil {
				return containerIDs, errors.Wrap(err, "error recording sidecar service container")
			}

			/**
			 * We need to record the logs for service containers in the background, this is not ideal BUT
			 * it does mean we can have real time logs when running them.
			 * TODO handle this better and capture error?
			 */
			stopWaiting := make(chan bool)
			timeoutDuration := defaultWaitTimeout
			if sidecar.Wait != nil && sidecar.Wait.Timeout != nil {
				timeoutDuration = time.Duration(*sidecar.Wait.Timeout) * time.Second
			}
			timeout := time.NewTimer(timeoutDuration)

			var regex *regexp.Regexp
			if sidecar.Wait != nil {
				log.Println("waiting for sidecar service container output to match regex: ", sidecar.Wait.Output)
				regex = regexp.MustCompile(sidecar.Wait.Output)

				go func() {
					<-timeout.C
					stopWaiting <- true
				}()
			}

			go func() {
				_ = pipeline.Runtime.CopyLogsForContainer(
					context,
					containerID,
					&util.LoggerWriter{
						Recorder: func(logLine string) error {
							if regex != nil && regex.MatchString(logLine) {
								stopWaiting <- true
							}
							return pipeline.Recorder.RecordContainerLog(containerID, logLine, shared.LogTypeStdOut)
						},
					},
					&util.LoggerWriter{
						Recorder: func(logLine string) error {
							if regex != nil && regex.MatchString(logLine) {
								stopWaiting <- true
							}
							return pipeline.Recorder.RecordContainerLog(containerID, logLine, shared.LogTypeStdErr)
						},
					},
				)
			}()

			if sidecar.Wait != nil {
				<-stopWaiting
				regex = nil
				if !timeout.Stop() {
					return containerIDs, errors.New(fmt.Sprintf("error waiting for sidecar service container output to match regex %s", sidecar.Wait.Output))
				}
			}
		}
	}

	// Now dispatch all of our step containers
	if stage.Steps != nil {
		for _, container := range stage.Steps {
			// First create the container and add the ID to our container IDs
			// We also add our container ID to our list of containers to terminate afterwards
			containerID, err := pipeline.Runtime.DispatchContainer(context, jobID, container)
			if containerID != shared.EmptyContainerID {
				containerIDs = append(containerIDs, containerID)
			}
			if err != nil {
				return containerIDs, errors.Wrap(err, "error dispatching step container")
			}

			if err = pipeline.Recorder.RecordContainer(jobID, containerID, shared.ContainerMeta{StageID: stageID, Service: false}, container, shared.ContainerStateStarting); err != nil {
				return containerIDs, errors.Wrap(err, "error recording step container creation")
			}

			// We want our container to be running or stopped (stopped is ok if the command execs really quickly)
			if e := pipeline.Runtime.WaitForContainer(
				context,
				containerID,
				shared.ContainerWaitCondition{State: shared.ContainerWaitStopped | shared.ContainerWaitRunning},
			); e != nil {
				err = errors.Wrap(e, "error waiting for container to be ready")
			}

			// FindAllByJobID the logs from the container, THIS WILL BLOCK until the container stops, i.e it runs to completion
			if e := pipeline.Runtime.CopyLogsForContainer(
				context,
				containerID,
				&util.LoggerWriter{
					Recorder: func(log string) error {
						return pipeline.Recorder.RecordContainerLog(containerID, log, shared.LogTypeStdOut)
					},
				},
				&util.LoggerWriter{
					Recorder: func(log string) error {
						return pipeline.Recorder.RecordContainerLog(containerID, log, shared.LogTypeStdErr)
					},
				},
			); e != nil {
				err = util.ErrorAppend(errors.Wrap(e, "error copying container logs"), err)
			}

			// If we made it this far, lets terminate the container as we are done with this step
			if e := pipeline.Runtime.TerminateContainer(context, containerID); e != nil {
				err = util.ErrorAppend(errors.Wrap(e, "error terminating step container"), err)
			}

			// Remove the container id we just terminated from our slice of container ids
			containerIDs = containerIDs[:len(containerIDs)-1]

			// Now record the container state
			containerState := shared.ContainerStateStopped
			if err != nil {
				containerState = shared.ContainerStateError
			}
			if e := pipeline.Recorder.RecordContainerState(containerID, containerState); e != nil {
				err = util.ErrorAppend(errors.Wrap(e, "error recording step container state"), err)
			}

			if err != nil {
				return containerIDs, util.ErrorAppend(errors.New("error executing container"), err)
			}
		}
	}

	return containerIDs, nil
}

func (pipeline *Pipeline) Execute(ctx context.Context, spec shared.Spec, workingDir string, jobID shared.JobID) error {

	for stageID, stage := range spec.Stages {

		log.Println("initializing stage runtime")
		err := pipeline.Runtime.Initialize(ctx, jobID, workingDir)
		if err != nil {
			err = errors.Wrap(err, "error initializing stage container runtime")
		}

		e := pipeline.Recorder.RecordStageState(jobID, stageID, shared.StageStateRunning)
		if e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, "error recording stage state"))
		}

		// Here we need to execute the stage and cleanup left over containers
		// If we get an error, dont return instead set the error and handle it at the end.
		// This way we can pass them back up the stack
		if err == nil {
			containerIds, e := pipeline.executeStage(ctx, jobID, stageID, stage)
			if e != nil {
				err = util.ErrorAppend(err, errors.Wrap(e, fmt.Sprintf("error running %s stage", stageID)))
			}

			e = pipeline.cleanUp(context.Background(), containerIds)
			if e != nil {
				err = util.ErrorAppend(err, errors.Wrap(e, "error cleaning up containers from stage"))
			}
		}

		state := shared.StageStateSuccess
		if err != nil {
			state = shared.StageStateError
		}
		e = pipeline.Recorder.RecordStageState(jobID, stageID, state)
		if e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, "error recording stage state"))
		}

		e = pipeline.Runtime.Terminate(context.Background(), jobID)
		if e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, "error terminating container"))
		}

		if err != nil {
			e = pipeline.Recorder.RecordLog(jobID, err.Error(), shared.LogTypeStdErr, stageID)
			return util.ErrorAppend(err, errors.Wrap(e, "error recording failure"))
		}
	}

	return nil
}

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
			return util.ErrorAppend(err, errors.Wrap(e, "error terminating container"))
		}

		if e := pipeline.Recorder.RecordContainerState(containerID, shared.ContainerStateStopped); e != nil {
			return util.ErrorAppend(err, errors.Wrap(e, "error updating terminating container status"))
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
				return containerIDs, errors.Wrap(err, "error dispatching step service container")
			}

			// Record our container as starting
			err = pipeline.Recorder.RecordContainer(jobID, containerID, shared.ContainerMeta{StageID: stageID, Service: true}, sidecar, shared.ContainerStateStarting)
			if err != nil {
				return containerIDs, errors.Wrap(err, "error recording step service container creation")
			}

			// Wait for our container to be stopWaiting, then mark it as running
			if err = pipeline.Runtime.WaitForContainer(context, containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitRunning}); err != nil {
				return containerIDs, errors.Wrap(err, "error waiting for step container to be stopWaiting")
			}

			if err = pipeline.Recorder.RecordContainerState(containerID, shared.ContainerStateRunning); err != nil {
				return containerIDs, errors.Wrap(err, "error recording step container state")
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
				log.Println("waiting for container output to match regex: ", sidecar.Wait.Output)
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
					return containerIDs, errors.New(fmt.Sprintf("error waiting for container output to match regex %s", sidecar.Wait.Output))
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
			if err = pipeline.Runtime.WaitForContainer(
				context,
				containerID,
				shared.ContainerWaitCondition{State: shared.ContainerWaitStopped | shared.ContainerWaitRunning},
			); err != nil {
				return containerIDs, errors.Wrap(err, "error waiting for step container to be ready")
			}

			if err = pipeline.Recorder.RecordContainerState(containerID, shared.ContainerStateRunning); err != nil {
				return containerIDs, errors.Wrap(err, "error recording step container state")
			}

			// FindAllByJobID the logs from the container, THIS WILL BLOCK until the container stops, i.e it runs to completion
			err = pipeline.Runtime.CopyLogsForContainer(
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
			)
			if err != nil {
				return containerIDs, errors.Wrap(err, "error copying step container logs")
			}

			// If we made it this far, lets terminate the container as we are done with this step
			if err = pipeline.Runtime.TerminateContainer(context, containerID); err != nil {
				return containerIDs, errors.Wrap(err, "error terminating step container")
			}

			// Remove the container id we just terminated from our slice of container ids
			containerIDs = containerIDs[:len(containerIDs)-1]

			if err = pipeline.Recorder.RecordContainerState(containerID, shared.ContainerStateStopped); err != nil {
				return containerIDs, errors.Wrap(err, "error recording step container state")
			}
		}
	}

	return containerIDs, nil
}

func (pipeline *Pipeline) Execute(ctx context.Context, spec shared.Spec, workingDir string, jobID shared.JobID) error {
	log.Println("initializing runtime")
	err := pipeline.Runtime.Initialize(ctx, jobID, workingDir)
	if err != nil {
		return errors.Wrap(err, "error initializing runtime config")
	}

	for stageID, stage := range spec.Stages {
		e := pipeline.Recorder.RecordStageState(jobID, stageID, shared.StageStateRunning)
		if e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, "error recording stage state"))
			break
		}

		// Here we need to execute the stage and cleanup left over containers
		// If we get an error, dont return instead set the error and handle it at the end.
		// This way we can pass them back up the stack
		containerIds, e := pipeline.executeStage(ctx, jobID, stageID, stage)
		if e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, fmt.Sprintf("error running %s stage", stageID)))
		}

		e = pipeline.cleanUp(context.Background(), containerIds)
		if e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, "error cleaning up containers from stage"))
		}

		state := shared.StageStateSuccess
		if err != nil {
			state = shared.StageStateError
		}
		e = pipeline.Recorder.RecordStageState(jobID, stageID, state)
		if e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, "error recording stage state"))
			break
		}
	}

	e := pipeline.Runtime.Terminate(context.Background(), jobID)
	if e != nil {
		err = util.ErrorAppend(err, errors.Wrap(e, "error terminating container"))
	}

	return err
}

/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package pipeline

import (
	"context"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/recorder"
	"go-brunel/internal/pkg/runner/runtime"
	"go-brunel/internal/pkg/runner/trigger"
	"go-brunel/internal/pkg/shared"
	"go-brunel/internal/pkg/shared/util"
	"log"
)

type JobHandler struct {
	RuntimeFactory runtime.Factory
	Recorder       recorder.Recorder
	WorkSpace      WorkSpace
}

const (
	endStage = "end"
)

// Handle will process a job trigger event and will record the status of the job
func (handler *JobHandler) Handle(event trigger.Event) {
	log.Printf("running in directory: %s\n", event.WorkDir)
	if err := handler.processJob(event.Context, event); err != nil {
		if err := handler.Recorder.RecordLog(event.Job.ID, err.Error(), shared.LogTypeStdErr, endStage); err != nil {
			log.Println(err)
		}
		if event.Context.Err() != nil {
			event.Job.State = shared.JobStateCancelled
		} else {
			event.Job.State = shared.JobStateFailed
		}
	} else {
		event.Job.State = shared.JobStateSuccess
	}
	event.JobState <- event.Job.State
}

// processJob should execute the full pipeline returning any errors to our caller
func (handler *JobHandler) processJob(context context.Context, event trigger.Event) error {
	pipelineRuntime, err := handler.RuntimeFactory.Create()
	if err != nil {
		return errors.Wrap(err, "error creating pipeline runtime")
	}

	pipelineSpec, err := handler.WorkSpace.Prepare(event)
	if err != nil {
		return util.ErrorAppend(
			errors.Wrap(err, "failed to prepare workspace"),
			handler.WorkSpace.CleanUp(event),
		)
	}

	pipeline := Pipeline{
		Runtime:  pipelineRuntime,
		Recorder: handler.Recorder,
	}
	err = errors.Wrap(
		pipeline.Execute(context, *pipelineSpec, event.WorkDir, event.Job.ID),
		"failed to execute pipeline",
	)
	return util.ErrorAppend(
		err,
		handler.WorkSpace.CleanUp(event),
	)
}

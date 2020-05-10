package pipeline

import (
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/environment"
	"go-brunel/internal/pkg/runner/parser"
	"go-brunel/internal/pkg/runner/recorder"
	"go-brunel/internal/pkg/runner/trigger"
	"go-brunel/internal/pkg/runner/vcs"
	"go-brunel/internal/pkg/shared"
	"go-brunel/internal/pkg/shared/util"
	"os"
)

type WorkSpace interface {
	Prepare(event trigger.Event) (*shared.Spec, error)
	CleanUp(event trigger.Event) error
}

type LocalWorkSpace struct {
	VCS                vcs.VCS
	EnvironmentFactory environment.Factory
	Recorder           recorder.Recorder
}

const (
	pipelineFile = ".brunel.jsonnet"

	preparingStageID  shared.StageID = "prepare"
	cleaningUpStageID shared.StageID = "clean"
)

func (w *LocalWorkSpace) Prepare(event trigger.Event) (*shared.Spec, error) {
	err := w.Recorder.RecordStageState(event.Job.ID, preparingStageID, shared.StageStateRunning)
	if err != nil {
		return nil, errors.Wrap(
			err,
			"error recording job state",
		)
	}

	_ = w.Recorder.RecordLog(event.Job.ID, "preparing workspace", shared.LogTypeStdOut, preparingStageID)

	progress := &util.LoggerWriter{
		Recorder: func(log string) error {
			return w.Recorder.RecordLog(event.Job.ID, log, shared.LogTypeStdOut, preparingStageID)
		},
	}

	if event.Job.Repository != trigger.LocalRepository {
		if err := w.VCS.Clone(vcs.Options{
			Directory:     event.WorkDir,
			RepositoryURL: event.Job.Repository.URI,
			Branch:        event.Job.Commit.Branch,
			Revision:      event.Job.Commit.Revision,
			Progress:      progress,
		}); err != nil {
			e := w.Recorder.RecordStageState(event.Job.ID, preparingStageID, shared.StageStateError)

			return nil, errors.Wrap(
				util.ErrorAppend(err, util.ErrorAppend(progress.Close(), e)),
				"error cloning repository",
			)
		}
	}

	_ = w.Recorder.RecordLog(event.Job.ID, "parsing specification", shared.LogTypeStdOut, preparingStageID)

	p := parser.JsonnetParser{
		Event: event,
		VCS:   w.VCS,
		EnvironmentProvider: w.EnvironmentFactory.Create(
			event.Job.EnvironmentID,
		),
	}

	spec, err := p.Parse(pipelineFile, progress)
	if err != nil {
		err = errors.Wrap(err, "error parsing pipeline specification")
	}

	if e := progress.Close(); e != nil {
		err = util.ErrorAppend(err, e)
	}

	stageState := shared.StageStateSuccess
	if err != nil {
		stageState = shared.StageStateError
	}

	e := w.Recorder.RecordStageState(event.Job.ID, preparingStageID, stageState)
	if e != nil {
		err = util.ErrorAppend(err, e)
	}

	return spec, err
}

func (w *LocalWorkSpace) CleanUp(event trigger.Event) error {
	err := w.Recorder.RecordStageState(event.Job.ID, cleaningUpStageID, shared.StageStateRunning)
	if err != nil {
		err = errors.Wrap(err, "error recording stage")
	}

	_ = w.Recorder.RecordLog(event.Job.ID, "cleaning workspace", shared.LogTypeStdOut, cleaningUpStageID)

	if event.Job.Repository != trigger.LocalRepository {
		if e := os.RemoveAll(event.WorkDir); e != nil {
			err = util.ErrorAppend(err, errors.Wrap(e, "error cleaning up workspace"))
		}
	}

	stageState := shared.StageStateSuccess
	if err != nil {
		stageState = shared.StageStateError
	}

	e := w.Recorder.RecordStageState(event.Job.ID, cleaningUpStageID, stageState)
	if e != nil {
		err = util.ErrorAppend(err, errors.Wrap(e, "error recording stage"))
	}

	return nil
}

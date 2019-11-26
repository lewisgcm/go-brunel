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
	pipelineFile   = ".brunel.jsonnet"
	workspaceStage = "preparing"
)

func (w *LocalWorkSpace) Prepare(event trigger.Event) (*shared.Spec, error) {
	progress := &util.LoggerWriter{
		Recorder: func(log string) error {
			return w.Recorder.RecordLog(event.Job.ID, log, shared.LogTypeStdOut, workspaceStage)
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
			return nil, errors.Wrap(
				util.ErrorAppend(err, progress.Close()),
				"error cloning repository",
			)
		}
	}

	p := parser.JsonnetParser{
		WorkingDirectory: event.WorkDir,
		VCS:              w.VCS,
		EnvironmentProvider: w.EnvironmentFactory.Create(
			[]string{event.Job.Repository.Project, event.Job.Repository.Name},
		),
	}

	spec, err := p.Parse(pipelineFile, progress)
	if err != nil {
		err = errors.Wrap(err, "error parsing pipeline specification")
	}

	return spec, util.ErrorAppend(err, progress.Close())
}

func (w *LocalWorkSpace) CleanUp(event trigger.Event) error {
	if event.Job.Repository != trigger.LocalRepository {
		if e := os.RemoveAll(event.WorkDir); e != nil {
			return errors.Wrap(e, "error cleaning up workspace")
		}
	}
	return nil
}

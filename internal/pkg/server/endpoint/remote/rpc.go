package remote

import (
	"go-brunel/internal/pkg/server/notify"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"go-brunel/internal/pkg/shared/remote"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

type RPC struct {
	Notify          notify.Notify
	JobStore        store.JobStore
	LogStore        store.LogStore
	ContainerStore  store.ContainerStore
	RepositoryStore store.RepositoryStore
	StageStore      store.StageStore
}

func (t *RPC) GetNextAvailableJob(args *remote.Empty, reply *remote.GetNextAvailableJobResponse) error {
	job, e := t.JobStore.Next()
	if e != nil {
		return errors.Wrap(e, "error getting next job from store")
	}

	if job != nil {
		r, e := t.RepositoryStore.Get(job.RepositoryID)
		if e != nil {
			return errors.Wrap(e, "error getting job repository from store")
		}
		log.Info("job with id ", job.ID, " has started")

		reply.Job = &shared.Job{
			ID:     job.ID,
			State:  job.State,
			Commit: job.Commit,
			Repository: shared.Repository{
				URI:     r.URI,
				Name:    r.Name,
				Project: r.Project,
			},
		}
	}
	return nil
}

func (t *RPC) SetJobState(args *remote.SetJobStateRequest, reply *remote.Empty) error {
	if args.State > shared.JobStateProcessing {
		log.Info("job with id ", args.Id, " has stopped")

		if err := t.JobStore.UpdateStoppedAtByID(args.Id, time.Now()); err != nil {
			return errors.Wrap(err, "error storing job stop time")
		}
	}

	if err := t.JobStore.UpdateStateByID(args.Id, args.State); err != nil {
		return errors.Wrap(err, "error storing job state")
	}

	return errors.Wrap(
		t.Notify.Notify(args.Id),
		"error notifying job status",
	)
}

func (t *RPC) HasBeenCancelled(args *shared.JobID, reply *bool) error {
	c, e := t.JobStore.Get(*args)
	*reply = c.State == shared.JobStateCancelled
	return e
}

func (t *RPC) Log(args *remote.LogRequest, reply *remote.Empty) error {
	return errors.Wrap(
		t.LogStore.Log(store.Log{
			JobID:   args.Id,
			Message: args.Message,
			LogType: args.LogType,
			StageID: args.StageID,
			Time:    time.Now(),
		}),
		"error storing log",
	)
}

func (t *RPC) SetStageState(args *remote.SetStageStateRequest, reply *remote.Empty) error {
	var startTime time.Time
	var stopTime time.Time

	if args.State > shared.StageStateRunning {
		stopTime = time.Now()
	} else {
		startTime = time.Now()
	}

	err := t.StageStore.AddOrUpdate(store.Stage{
		ID:        args.Id,
		JobID:     args.JobID,
		State:     args.State,
		StartedAt: &startTime,
		StoppedAt: &stopTime,
	})

	return errors.Wrap(err, "error storing stage stopped time")
}

func (t *RPC) AddContainer(args *remote.AddContainerRequest, reply *remote.Empty) error {
	return errors.Wrap(
		t.ContainerStore.Add(store.Container{
			JobID:       args.Id,
			ContainerID: args.ContainerID,
			Meta:        args.Meta,
			Spec:        args.Container,
			State:       args.State,
			CreatedAt:   time.Now(),
		}),
		"error storing container",
	)
}

func (t *RPC) SetContainerState(args *remote.SetContainerStateRequest, reply *remote.Empty) error {
	if args.State == shared.ContainerStateStopped {
		if err := t.ContainerStore.UpdateStoppedAtByContainerID(args.Id, time.Now()); err != nil {
			return errors.Wrap(err, "error storing container stop time")
		}
	} else if args.State == shared.ContainerStateRunning || args.State == shared.ContainerStateStarting {
		if err := t.ContainerStore.UpdateStartedAtByContainerID(args.Id, time.Now()); err != nil {
			return errors.Wrap(err, "error storing container start time")
		}
	}
	return errors.Wrap(
		t.ContainerStore.UpdateStateByContainerID(args.Id, args.State),
		"error storing container state",
	)
}

func (t *RPC) ContainerLog(args *remote.ContainerLogRequest, reply *remote.Empty) error {
	return errors.Wrap(
		t.LogStore.ContainerLog(store.ContainerLog{
			ContainerID: args.Id,
			Message:     args.Message,
			LogType:     args.LogType,
			Time:        time.Now(),
		}),
		"error storing container log",
	)
}

func (t *RPC) SearchForSecret(args *remote.SearchForXRequest, reply *string) error {
	*reply = ""
	return nil
}

func (t *RPC) SearchForValue(args *remote.SearchForXRequest, reply *string) error {
	*reply = ""
	return nil
}

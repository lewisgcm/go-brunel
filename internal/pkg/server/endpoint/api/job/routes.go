package job

import (
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/security"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
)

type jobHandler struct {
	jobStore        store.JobStore
	logStore        store.LogStore
	stageStore      store.StageStore
	containerStore  store.ContainerStore
	repositoryStore store.RepositoryStore
	jwtSerializer   security.TokenSerializer
}

func (handler *jobHandler) get(r *http.Request) api.Response {
	id := chi.URLParam(r, "id")
	job, err := handler.jobStore.Get(shared.JobID(id))
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(err, "error getting job")
	}

	repository, err := handler.repositoryStore.Get(job.RepositoryID)
	if err != nil {
		return api.InternalServerError(err, "error getting job")
	}

	return api.Ok(struct {
		store.Job
		Repository store.Repository
	}{
		Job:        job,
		Repository: repository,
	})
}

func (handler *jobHandler) progress(r *http.Request) api.Response {
	id := shared.JobID(chi.URLParam(r, "id"))
	since, err := api.ParseQueryTime(r, "since", false, time.Time{})
	if err != nil {
		return api.InternalServerError(err, "error parsing query parameter 'since'")
	}

	details := struct {
		State  shared.JobState
		Stages []struct {
			store.Stage
			Containers []store.Container
			Logs       []store.Log
		}
	}{}

	job, err := handler.jobStore.Get(id)
	if err != nil {
		return api.InternalServerError(err, "error getting job")
	}
	details.State = job.State

	stages, err := handler.stageStore.FindAllByJobID(id)
	if err != nil {
		return api.InternalServerError(err, "error getting job stages")
	}

	// Read out containers with a matching job id
	containers, err := handler.containerStore.FilterByJobID(id)
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(err, "error getting job containers")
	}

	// Read our the job level logs with the
	logs, err := handler.logStore.FilterLogByJobIDFromTime(id, since)
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(err, "error getting job logs")
	}

	// Map out out object for reading the UI
	for _, stage := range stages {

		var stageContainers []store.Container

		for _, c := range containers {
			if c.Meta.StageID != stage.ID {
				continue
			}

			// Overwrite the environment, we dont want to leak any sensitive information
			c.Spec.Environment = nil

			stageContainers = append(stageContainers, c)
		}

		var stageLogs []store.Log
		for _, l := range logs {
			if l.StageID == stage.ID {
				stageLogs = append(stageLogs, l)
			}
		}

		mappedStage := struct {
			store.Stage
			Containers []store.Container
			Logs       []store.Log
		}{
			Stage:      stage,
			Containers: stageContainers,
			Logs:       stageLogs,
		}

		details.Stages = append(details.Stages, mappedStage)
	}

	return api.Ok(details)
}

func (handler *jobHandler) cancel(r *http.Request) api.Response {
	id := chi.URLParam(r, "id")
	identity, err := handler.jwtSerializer.Decode(r)
	if err != nil {
		return api.InternalServerError(err, "error decoding token")
	}

	job, err := handler.jobStore.Get(shared.JobID(id))
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(err, "error getting job")
	}

	if !(identity.Username != job.StartedBy || identity.Role != security.UserRoleAdmin) {
		return api.UnAuthorized()
	}

	err = handler.jobStore.CancelByID(shared.JobID(id), identity.Username)
	if err != nil {
		return api.InternalServerError(err, "error setting job state")
	}

	log.Info("job with id ", id, " has been cancelled")
	return api.NoContent()
}

func Routes(
	jobStore store.JobStore,
	logStore store.LogStore,
	stageStore store.StageStore,
	containerStore store.ContainerStore,
	repositoryStore store.RepositoryStore,
	jwtSerializer security.TokenSerializer,
) *chi.Mux {
	handler := jobHandler{
		jobStore:        jobStore,
		logStore:        logStore,
		stageStore:      stageStore,
		repositoryStore: repositoryStore,
		containerStore:  containerStore,
		jwtSerializer:   jwtSerializer,
	}
	router := chi.NewRouter()
	router.Get("/{id}", api.Handle(handler.get))
	router.Get("/{id}/progress", api.Handle(handler.progress))
	router.Delete("/{id}", api.Handle(handler.cancel))
	return router
}

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
	jobStore       store.JobStore
	logStore       store.LogStore
	stageStore     store.StageStore
	containerStore store.ContainerStore
	jwtSerializer  security.TokenSerializer
}

func (handler *jobHandler) get(r *http.Request) (interface{}, int, error) {
	id := chi.URLParam(r, "id")
	job, err := handler.jobStore.Get(shared.JobID(id))
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(err, "error getting job")
	}

	return job, http.StatusOK, nil
}

func (handler *jobHandler) progress(r *http.Request) (interface{}, int, error) {
	id := shared.JobID(chi.URLParam(r, "id"))
	since, err := api.ParseQueryTime(r, "since", false, time.Time{})
	if err != nil {
		return api.InternalServerError(err, "error parsing query parameter 'since'")
	}

	details := struct {
		Job    store.Job
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
	details.Job = job

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

	return details, http.StatusOK, nil
}

func (handler *jobHandler) cancel(r *http.Request) (interface{}, int, error) {
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

	err = handler.jobStore.UpdateStateByID(shared.JobID(id), shared.JobStateCancelled)
	if err != nil {
		return api.InternalServerError(err, "error setting job state")
	}

	log.Info("job with id ", id, " has been cancelled")
	return nil, http.StatusNoContent, nil
}

func Routes(
	jobStore store.JobStore,
	logStore store.LogStore,
	containerStore store.ContainerStore,
	stageStore store.StageStore,
	jwtSerializer security.TokenSerializer,
) *chi.Mux {
	handler := jobHandler{
		jobStore:       jobStore,
		logStore:       logStore,
		stageStore:     stageStore,
		containerStore: containerStore,
		jwtSerializer:  jwtSerializer,
	}
	router := chi.NewRouter()
	router.Get("/{id}", api.Handle(handler.get))
	router.Get("/{id}/progress", api.Handle(handler.progress))
	router.Delete("/{id}", api.Handle(handler.cancel))
	return router
}

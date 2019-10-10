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
		Containers    []store.Container
		Logs          []store.Log
		ContainerLogs []store.ContainerLog
	}{}

	// Read out containers with a matching job id
	c, err := handler.containerStore.FilterByCreatedTimeAndJobID(id, since)
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(err, "error getting job containers")
	}
	details.Containers = c

	// Read our the job level logs with the
	l, err := handler.logStore.FilterLogByJobIDFromTime(id, since)
	if err != nil {
		if err == store.ErrorNotFound {
			return api.NotFound()
		}
		return api.InternalServerError(err, "error getting job logs")
	}
	details.Logs = l

	// Read out the container level logs
	for _, c := range details.Containers {
		l, err := handler.logStore.FilterContainerLogByContainerIDFromTime(c.ContainerID, since)
		if err != nil {
			return api.InternalServerError(err, "error container logs")
		}
		details.ContainerLogs = append(details.ContainerLogs, l...)
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
	jwtSerializer security.TokenSerializer,
) *chi.Mux {
	handler := jobHandler{
		jobStore:       jobStore,
		logStore:       logStore,
		containerStore: containerStore,
		jwtSerializer:  jwtSerializer,
	}
	router := chi.NewRouter()
	router.Get("/{id}", api.Handle(handler.get))
	router.Get("/{id}/progress", api.Handle(handler.progress))
	router.Delete("/{id}", api.Handle(handler.cancel))
	return router
}

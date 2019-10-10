package container

import (
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/security"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type jobHandler struct {
	logStore      store.LogStore
	jwtSerializer security.TokenSerializer
}

func (handler *jobHandler) logs(r *http.Request) (interface{}, int, error) {
	id := shared.ContainerID(chi.URLParam(r, "id"))
	since, err := api.ParseQueryTime(r, "since", false, time.Time{})
	if err != nil {
		return api.BadRequest(err, "invalid 'since' query parameter")
	}

	logs, err := handler.logStore.FilterContainerLogByContainerIDFromTime(id, since)
	if err != nil {
		return api.InternalServerError(err, "error getting container logs")
	}
	return logs, http.StatusOK, nil
}

func Routes(repository store.LogStore, jwtSerializer security.TokenSerializer) *chi.Mux {
	handler := jobHandler{
		logStore:      repository,
		jwtSerializer: jwtSerializer,
	}
	router := chi.NewRouter()
	router.Get("/{id}/logs", api.Handle(handler.logs))
	return router
}

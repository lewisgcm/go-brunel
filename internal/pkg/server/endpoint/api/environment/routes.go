package environment

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"net/http"
)

type environmentHandler struct {
	environmentStore store.EnvironmentStore
}

func (handler *environmentHandler) list(r *http.Request) api.Response {
	filter := r.URL.Query().Get("filter")
	entities, err := handler.environmentStore.Filter(filter)
	if err != nil {
		return api.InternalServerError(errors.Wrap(err, "internal error"))
	}
	return api.Ok(entities)
}

func (handler *environmentHandler) get(r *http.Request) api.Response {
	id := chi.URLParam(r, "id")
	environment, err := handler.environmentStore.Get(shared.EnvironmentID(id))
	if err != nil {
		return api.InternalServerError(errors.Wrap(err, "error getting environment"))
	}
	return api.Ok(environment)
}

func (handler *environmentHandler) save(r *http.Request) api.Response {
	environment := store.Environment{}
	if e := json.NewDecoder(r.Body).Decode(&environment); e != nil {
		return api.BadRequest(e, "bad request")
	}

	environment.Clean()
	if e := environment.IsValid(); e != nil {
		return api.BadRequest(errors.Wrap(e, "invalid environment"), e.Error())
	}

	result, err := handler.environmentStore.AddOrUpdate(environment)
	if err != nil {
		return api.InternalServerError(errors.Wrap(err, "error saving environment"))
	}
	return api.Ok(result)
}

func Routes(environmentStore store.EnvironmentStore) *chi.Mux {
	handler := environmentHandler{
		environmentStore: environmentStore,
	}
	router := chi.NewRouter()
	router.Get("/", api.Handle(handler.list))
	router.Get("/{id}", api.Handle(handler.get))
	router.Post("/", api.Handle(handler.save))

	return router
}

package repository

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/store"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	defaultJobPageSize = 5
	maxJobPageSize     = 20
)

type repositoryHandler struct {
	jobStore        store.JobStore
	repositoryStore store.RepositoryStore
}

func (handler *repositoryHandler) jobs(r *http.Request) api.Response {
	id := chi.URLParam(r, "id")
	filter := r.URL.Query().Get("filter")
	sortColumn := r.URL.Query().Get("sortColumn")
	sortOrder := -1

	if r.URL.Query().Get("sortOrder") == "asc" {
		sortOrder = 1
	}

	if sortColumn != "create_at" && sortColumn != "state" {
		sortColumn = "created_at"
	}

	pageIndex, err := api.ParseQueryInt(r, "pageIndex", false, 0)
	if err != nil {
		return api.BadRequest(err, "error parsing pageIndex query parameter")
	}

	pageSize, err := api.ParseQueryInt(r, "pageSize", false, defaultJobPageSize)
	if err != nil {
		return api.BadRequest(err, "error parsing pageSize query parameter")
	}

	if pageSize > maxJobPageSize {
		return api.BadRequest(err, "requested page size is above limit")
	}

	jobs, err := handler.jobStore.FilterByRepositoryID(id, filter, pageIndex, pageSize, string(sortColumn), sortOrder)
	if err != nil {
		return api.InternalServerError(err, "error getting repository jobs")
	}
	return api.Ok(jobs)
}

func (handler *repositoryHandler) list(r *http.Request) api.Response {
	filter := r.URL.Query().Get("filter")
	repositories, err := handler.repositoryStore.Filter(filter)
	if err != nil {
		return api.InternalServerError(err, "error getting repository details")
	}
	return api.Ok(repositories)
}

func (handler *repositoryHandler) get(r *http.Request) api.Response {
	id := chi.URLParam(r, "id")
	repositories, err := handler.repositoryStore.Get(store.RepositoryID(id))
	if err != nil {
		return api.InternalServerError(err, "error getting repository")
	}
	return api.Ok(repositories)
}

func (handler *repositoryHandler) setTriggers(r *http.Request) api.Response {
	id := chi.URLParam(r, "id")
	triggers := []store.RepositoryTrigger{}
	err := json.NewDecoder(r.Body).Decode(&triggers)
	if err != nil {
		return api.InternalServerError(
			errors.Wrap(err, "error decoding json"),
			"internal error",
		)
	}

	for _, t := range triggers {
		if !t.IsValid() {
			return api.BadRequest(errors.New("invalid trigger supplied"), "invalid trigger")
		}
	}

	err = handler.repositoryStore.SetTriggers(store.RepositoryID(id), triggers)
	if err != nil {
		return api.InternalServerError(
			errors.Wrap(err, "error saving environment"),
			"internal error",
		)
	}

	return api.NoContent()
}

func Routes(repositoryStore store.RepositoryStore, jobStore store.JobStore) *chi.Mux {
	handler := repositoryHandler{
		jobStore:        jobStore,
		repositoryStore: repositoryStore,
	}
	router := chi.NewRouter()
	router.Get("/", api.Handle(handler.list))
	router.Get("/{id}", api.Handle(handler.get))
	router.Get("/{id}/jobs", api.Handle(handler.jobs))
	router.Put("/{id}/triggers", api.Handle(handler.setTriggers))
	return router
}

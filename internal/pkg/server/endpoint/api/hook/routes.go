package hook

import (
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/notify"
	"go-brunel/internal/pkg/server/store"

	"github.com/go-chi/chi"
)

type webHookHandler struct {
	notifier        notify.Notify
	jobStore        store.JobStore
	repositoryStore store.RepositoryStore
}

func Routes(
	jobStore store.JobStore,
	repositoryStore store.RepositoryStore,
	notifier notify.Notify,
) *chi.Mux {
	handler := webHookHandler{
		jobStore:        jobStore,
		repositoryStore: repositoryStore,
		notifier:        notifier,
	}
	router := chi.NewRouter()
	router.Post("/gitlab", api.Handle(handler.gitLab))
	return router
}

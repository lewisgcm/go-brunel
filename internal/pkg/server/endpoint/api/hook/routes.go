package hook

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server"
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/notify"
	"go-brunel/internal/pkg/server/store"
	"regexp"

	"github.com/go-chi/chi"
)

type webHookHandler struct {
	configuration   server.WebHookConfiguration
	notifier        notify.Notify
	jobStore        store.JobStore
	repositoryStore store.RepositoryStore
}

func (handler *webHookHandler) finishHandling(repository store.Repository, job store.Job) api.Response {
	if !repository.IsValid() {
		return api.BadRequest(nil, "invalid project name or namespace supplied")
	}

	if !job.IsValid() {
		return api.BadRequest(nil, "invalid branch or revision supplied")
	}

	repo, err := handler.repositoryStore.AddOrUpdate(repository)
	if err != nil {
		return api.InternalServerError(errors.Wrap(err, "error storing github hook event repository"))
	}

	for _, t := range repo.Triggers {
		r, e := regexp.Compile(t.Pattern)
		if e != nil {
			return api.InternalServerError(
				errors.Wrap(e, "invalid pattern"),
			)
		}

		if r.Match([]byte(job.Commit.Branch)) {
			j, err := handler.jobStore.Add(store.Job{
				RepositoryID:  repo.ID,
				EnvironmentID: t.EnvironmentID,
				Commit:        job.Commit,
				State:         job.State,
				StartedBy:     job.StartedBy,
				CreatedAt:     job.CreatedAt,
			})
			if err != nil {
				return api.InternalServerError(errors.Wrap(err, "error storing hook event job"))
			}

			if err := handler.notifier.Notify(j.ID); err != nil {
				return api.InternalServerError(errors.Wrap(err, "error notifying job status from hook event"))
			}
		}
	}

	log.Info("received build notification hook for project ", repo.Project, "/", repo.Name)
	return api.NoContent()
}

func Routes(
	configuration server.WebHookConfiguration,
	jobStore store.JobStore,
	repositoryStore store.RepositoryStore,
	notifier notify.Notify,
) *chi.Mux {
	handler := webHookHandler{
		configuration:   configuration,
		jobStore:        jobStore,
		repositoryStore: repositoryStore,
		notifier:        notifier,
	}
	router := chi.NewRouter()
	router.Post("/gitlab", api.Handle(handler.gitLab))
	router.Post("/github", api.Handle(handler.gitHub))
	return router
}

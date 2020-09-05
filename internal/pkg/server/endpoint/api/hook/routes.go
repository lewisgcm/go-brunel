package hook

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server"
	"go-brunel/internal/pkg/server/bus"
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"regexp"
	"time"

	"github.com/go-chi/chi"
)

type webHookHandler struct {
	configuration   server.WebHookConfiguration
	jobStore        store.JobStore
	repositoryStore store.RepositoryStore
	bus             bus.EventBus
}

func (handler *webHookHandler) finishHandling(repository store.Repository, job store.Job) api.Response {
	job.Clean()
	repository.Clean()

	if e := repository.IsValid(); e != nil {
		return api.BadRequest(errors.Wrap(e, "invalid repository"), e.Error())
	}

	if e := job.IsValid(); e != nil {
		return api.BadRequest(errors.Wrap(e, "invalid job"), e.Error())
	}

	now := time.Now()
	repo, err := handler.repositoryStore.AddOrUpdate(repository)
	if err != nil {
		return api.InternalServerError(errors.Wrap(err, "error storing repository"))
	}

	if repo.CreatedAt.Unix() >= now.Unix() {
		if e := handler.bus.Send(shared.NewRepositoryCreated(repo.ID)); e != nil {
			return api.InternalServerError(errors.Wrap(err, "error trigger repository created event"))
		}
	}

	for _, t := range repo.Triggers {
		r, e := regexp.Compile(t.Pattern)
		if e != nil {
			return api.InternalServerError(
				errors.Wrap(e, "error compiling trigger pattern"),
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

			if e := handler.bus.Send(shared.NewJobCreated(j.ID, j.RepositoryID)); e != nil {
				return api.InternalServerError(errors.Wrap(e, "error create job event"))
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
	bus bus.EventBus,
) *chi.Mux {
	handler := webHookHandler{
		configuration:   configuration,
		jobStore:        jobStore,
		repositoryStore: repositoryStore,
		bus:             bus,
	}
	router := chi.NewRouter()
	router.Post("/gitlab", api.Handle(handler.gitLab))
	router.Post("/github", api.Handle(handler.gitHub))
	return router
}

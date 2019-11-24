package hook

import (
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

func (handler *webHookHandler) gitLab(r *http.Request) (interface{}, int, error) {
	hook, _ := gitlab.New()

	payload, err := hook.Parse(r, gitlab.PushEvents, gitlab.TagEvents)
	if err != nil && err != gitlab.ErrEventNotFound {
		return api.InternalServerError(err, "error handling gitlab hook event")
	}

	var job store.Job
	var repository store.Repository

	switch payload.(type) {
	case gitlab.PushEventPayload:
		event := payload.(gitlab.PushEventPayload)
		job = store.Job{
			Commit: shared.Commit{
				Branch:   event.Ref,
				Revision: event.After,
			},
			State:     shared.JobStateWaiting,
			StartedBy: event.UserEmail,
			CreatedAt: time.Now(),
		}
		repository = store.Repository{
			Project:   event.Project.Namespace,
			Name:      event.Project.Name,
			URI:       event.Project.GitHTTPURL,
			CreatedAt: time.Now(),
		}
	case gitlab.TagEventPayload:
		event := payload.(gitlab.PushEventPayload)
		job = store.Job{
			Commit: shared.Commit{
				Branch:   event.Ref,
				Revision: event.After,
			},
			State:     shared.JobStateWaiting,
			StartedBy: event.UserEmail,
			CreatedAt: time.Now(),
		}
		repository = store.Repository{
			Project:   event.Project.Namespace,
			Name:      event.Project.Name,
			URI:       event.Project.GitHTTPURL,
			CreatedAt: time.Now(),
		}
	default:
		return api.NotFound()
	}

	if !repository.IsValid() {
		return api.BadRequest(nil, "invalid project name or namespace supplied")
	}

	if !job.IsValid() {
		return api.BadRequest(nil, "invalid branch or revision supplied")
	}

	repo, err := handler.repositoryStore.AddOrUpdate(repository)
	if err != nil {
		return api.InternalServerError(err, "error storing gitlab hook event repository")
	}

	job.RepositoryID = repo.ID
	id, err := handler.jobStore.Add(job)
	if err != nil {
		return api.InternalServerError(err, "error storing gitlab hook event job")
	}

	if err := handler.notifier.Notify(id); err != nil {
		return api.InternalServerError(err, "error notifying job status from hook event")
	}

	log.Info("received build notification from gitlab hook for project ", repo.Project, "/", repo.Name)
	return nil, http.StatusOK, nil
}

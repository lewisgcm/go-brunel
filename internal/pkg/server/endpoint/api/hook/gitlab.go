package hook

import (
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"net/http"
	"time"

	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

func (handler *webHookHandler) gitLab(r *http.Request) api.Response {
	hook, _ := gitlab.New(gitlab.Options.Secret(handler.configuration.GitLabSecret))

	payload, err := hook.Parse(r, gitlab.PushEvents, gitlab.TagEvents)
	if err != nil && err != gitlab.ErrEventNotFound {
		return api.InternalServerError(errors.Wrap(err, "error handling gitlab hook event"))
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

	return handler.finishHandling(repository, job)
}

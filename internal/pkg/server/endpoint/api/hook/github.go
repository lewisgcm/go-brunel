package hook

import (
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"gopkg.in/go-playground/webhooks.v5/github"
	"net/http"
	"strings"
	"time"
)

func (handler *webHookHandler) gitHub(r *http.Request) api.Response {
	hook, _ := github.New(github.Options.Secret(handler.configuration.GitHubSecret))

	payload, err := hook.Parse(r, github.PushEvent)
	if err != nil && err != github.ErrEventNotFound {
		return api.InternalServerError(errors.Wrap(err, "error handling github hook event"))
	}

	var job store.Job
	var repository store.Repository

	switch payload.(type) {
	case github.PushPayload:
		event := payload.(github.PushPayload)
		job = store.Job{
			Commit: shared.Commit{
				Branch:   event.Ref,
				Revision: event.After,
			},
			State:     shared.JobStateWaiting,
			StartedBy: event.Pusher.Email,
			CreatedAt: time.Now(),
		}

		parts := strings.Split(event.Repository.FullName, "/")
		if len(parts) != 2 {
			return api.InternalServerError(errors.New("invalid project/repository name"))
		}

		repository = store.Repository{
			Project:   parts[0],
			Name:      parts[1],
			URI:       event.Repository.CloneURL,
			CreatedAt: time.Now(),
		}
	default:
		return api.NotFound()
	}

	return handler.finishHandling(repository, job)
}

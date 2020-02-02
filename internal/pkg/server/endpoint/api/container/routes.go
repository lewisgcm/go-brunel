package container

import (
	"github.com/buildkite/terminal-to-html"
	"go-brunel/internal/pkg/server/endpoint/api"
	"go-brunel/internal/pkg/server/security"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type jobHandler struct {
	logStore       store.LogStore
	containerStore store.ContainerStore
	jwtSerializer  security.TokenSerializer
}

func (handler *jobHandler) logs(w http.ResponseWriter, r *http.Request) {
	id := shared.ContainerID(chi.URLParam(r, "id"))
	since, err := api.ParseQueryTime(r, "since", false, time.Time{})
	if err != nil {
		log.Println("error parsing since time", err)
		return
	}

	state, err := handler.containerStore.GetContainerState(id)
	if err != nil {
		log.Println("error getting job state", err)
		return
	}

	logs, err := handler.logStore.FilterContainerLogByContainerIDFromTime(id, since)
	if err != nil {
		log.Println("error querying container logs", err)
		return
	}

	if *state == shared.ContainerStateStopped {
		w.Header().Add("X-Content-Complete", "True")
	}

	w.Header().Add("Content-Type", "text/html")
	for _, i := range logs {
		w.Write(terminal.Render([]byte(i.Message)))
		w.Write([]byte("<br/>"))
	}
}

func Routes(repository store.LogStore, containerStore store.ContainerStore, jwtSerializer security.TokenSerializer) *chi.Mux {
	handler := jobHandler{
		logStore:       repository,
		containerStore: containerStore,
		jwtSerializer:  jwtSerializer,
	}
	router := chi.NewRouter()
	router.Get("/{id}/logs", handler.logs)
	return router
}

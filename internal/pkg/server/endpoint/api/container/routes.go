package container

import (
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
	logStore      store.LogStore
	jwtSerializer security.TokenSerializer
}

func (handler *jobHandler) logs(w http.ResponseWriter,  r *http.Request) {
	id := shared.ContainerID(chi.URLParam(r, "id"))
	since, err := api.ParseQueryTime(r, "since", false, time.Time{})
	if err != nil {
		log.Println("error parsing since time", err)
		return
	}

	logs, err := handler.logStore.FilterContainerLogByContainerIDFromTime(id, since)
	if err != nil {
		log.Println("error querying container logs", err)
	}

	for _, i := range logs {
		w.Write([]byte(i.Message))
		w.Write([]byte("\n"))
	}
}

func Routes(repository store.LogStore, jwtSerializer security.TokenSerializer) *chi.Mux {
	handler := jobHandler{
		logStore:      repository,
		jwtSerializer: jwtSerializer,
	}
	router := chi.NewRouter()
	router.Get("/{id}/logs", handler.logs)
	return router
}

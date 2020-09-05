package event

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"go-brunel/internal/pkg/server/bus"
	"net/http"
)

type handler struct {
	bus bus.EventBus
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

const (
	channelBufferSize = 1000
)

func (h *handler) events(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("error upgrading web socket client:", err)
		return
	}

	events := make(chan interface{}, channelBufferSize)
	if e := h.bus.Listen(func(event interface{}) error {
		events <- event
		return nil
	}); e != nil {
		log.Errorf("error upgrading web socket client: ", e)
		return
	}

	defer c.Close()
	for {
		event := <-events
		if e := c.WriteJSON(event); e != nil {
			log.Errorf("error writing json to web socket client: %s", e)
			break
		}
	}
}

func Routes(
	bus bus.EventBus,
) *chi.Mux {
	h := handler{
		bus: bus,
	}

	router := chi.NewRouter()
	router.HandleFunc("/", h.events)

	return router
}

package notify

import (
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/bus"
	"go-brunel/internal/pkg/shared"
)

type listener struct {
	Bus    bus.EventBus
	Notify Notify
}

func NewListener(bus bus.EventBus, notify Notify) error {

	if e := bus.Listen(func(event interface{}) error {
		switch v := event.(type) {
		case shared.JobCreated:
			if e := notify.Notify(v.JobID); e != nil {
				return errors.Wrap(e, "error notifying job status")
			}
		}
		return nil
	}); e != nil {
		return errors.Wrap(e, "error listening for events")
	}

	return nil
}

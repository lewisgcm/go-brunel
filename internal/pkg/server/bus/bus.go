package bus

import "go-brunel/internal/pkg/shared/util"

type EventBus interface {
	Send(event interface{}) error

	Listen(listener func(event interface{}) error) error
}

type InMemoryEventBus struct {
	listeners []func(event interface{}) error
}

func (i *InMemoryEventBus) Send(event interface{}) error {
	var err error
	for _, l := range i.listeners {
		if e := l(event); e != nil {
			err = util.ErrorAppend(err, e)
		}
	}
	return err
}

func (i *InMemoryEventBus) Listen(listener func(event interface{}) error) error {
	i.listeners = append(i.listeners, listener)
	return nil
}

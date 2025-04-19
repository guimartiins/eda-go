package events

import "time"

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() any
}

type EventHandlerInterface interface {
	Handle(event EventInterface)
}

type EventDispatcherInterface interface {
	Subscribe(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Unsubscribe(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear() error
}

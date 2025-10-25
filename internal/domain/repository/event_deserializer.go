package repository

import (
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
)

type EventDeserializer interface {
	Deserialize(eventType string, eventData []byte) (event.Event, error)
}

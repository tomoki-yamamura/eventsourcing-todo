package deserializer

import (
	"fmt"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
)

type eventRegistry struct {
	deserializers map[string]eventDeserializer
}

type eventDeserializer interface {
	Deserialize(eventData []byte) (event.Event, error)
	EventType() string
}

func NewEventDeserializer() repository.EventDeserializer {
	registry := &eventRegistry{
		deserializers: make(map[string]eventDeserializer),
	}

	registry.register(NewTodoListCreatedEventDeserializer())
	registry.register(NewTodoAddedEventDeserializer())

	return registry
}

func (r *eventRegistry) register(deserializer eventDeserializer) {
	r.deserializers[deserializer.EventType()] = deserializer
}

func (r *eventRegistry) Deserialize(eventType string, eventData []byte) (event.Event, error) {
	deserializer, exists := r.deserializers[eventType]
	if !exists {
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	return deserializer.Deserialize(eventData)
}

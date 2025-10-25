package event

import (
	"time"

	"github.com/google/uuid"
)

type TodoListCreatedEvent struct {
	AggregateID uuid.UUID
	UserID      string
	EventID     uuid.UUID
	Timestamp   time.Time
	Version     int
}

func (e TodoListCreatedEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

func (e TodoListCreatedEvent) GetEventID() uuid.UUID {
	return e.EventID
}

func (e TodoListCreatedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e TodoListCreatedEvent) GetVersion() int {
	return e.Version
}

func (e TodoListCreatedEvent) GetEventType() string {
	return "TodoListCreatedEvent"
}
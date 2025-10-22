package event

import (
	"github.com/google/uuid"
	"time"
)

type Event interface {
	GetAggregateID() uuid.UUID
	GetEventID() uuid.UUID
	GetTimestamp() time.Time
	GetVersion() int
	GetEventType() string
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
)

type EventStore interface {
	SaveEvents(ctx context.Context, aggregateID uuid.UUID, events []event.Event) error
	LoadEvents(ctx context.Context, aggregateID uuid.UUID) ([]event.Event, error)
	GetAllEvents(ctx context.Context) ([]event.Event, error)
}

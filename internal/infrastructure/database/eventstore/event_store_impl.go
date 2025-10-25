package eventstore

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/transaction"
)

type eventStoreImpl struct{}

func NewEventStore() repository.EventStore {
	return &eventStoreImpl{}
}

func (e *eventStoreImpl) SaveEvents(ctx context.Context, aggregateID uuid.UUID, events []event.Event) error {
	tx, err := transaction.GetTx(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO events (
			aggregate_id, 
			event_id, 
			event_type, 
			event_data, 
			version, 
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, evt := range events {
		eventData, err := json.Marshal(evt)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query,
			aggregateID,
			evt.GetEventID(),
			evt.GetEventType(),
			eventData,
			evt.GetVersion(),
			time.Now(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *eventStoreImpl) LoadEvents(ctx context.Context, aggregateID uuid.UUID) ([]event.Event, error) {
	// TODO: Implement LoadEvents - event deserialization needed
	return nil, nil
}

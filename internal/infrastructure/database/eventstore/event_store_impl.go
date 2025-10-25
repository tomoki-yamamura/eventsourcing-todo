package eventstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/eventstore/deserializer"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/transaction"
)

type eventStoreImpl struct {
	deserializer repository.EventDeserializer
}

func NewEventStore() repository.EventStore {
	return &eventStoreImpl{
		deserializer: deserializer.NewEventDeserializer(),
	}
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
	tx, err := transaction.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT event_id, event_type, event_data, version, created_at
		FROM events 
		WHERE aggregate_id = $1 
		ORDER BY version ASC
	`

	rows, err := tx.QueryContext(ctx, query, aggregateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var eventID uuid.UUID
		var eventType string
		var eventData []byte
		var version int
		var createdAt time.Time

		err := rows.Scan(&eventID, &eventType, &eventData, &version, &createdAt)
		if err != nil {
			return nil, err
		}

		evt, err := e.deserializer.Deserialize(eventType, eventData)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize event %s: %w", eventType, err)
		}

		events = append(events, evt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

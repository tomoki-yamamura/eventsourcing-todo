package eventstore_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/config"
	domainevent "github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	appErrors "github.com/tomoki-yamamura/eventsourcing-todo/internal/errors"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/client"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/eventstore"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/transaction"
)

type testEvent struct {
	AggregateID uuid.UUID `json:"aggregate_id"`
	EventID     uuid.UUID `json:"event_id"`
	Type        string    `json:"event_type"`
	Version     int       `json:"version"`
	Title       string    `json:"title,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (e testEvent) GetAggregateID() uuid.UUID { return e.AggregateID }
func (e testEvent) GetEventID() uuid.UUID     { return e.EventID }
func (e testEvent) GetEventType() string      { return e.Type }
func (e testEvent) GetVersion() int           { return e.Version }
func (e testEvent) GetTimestamp() time.Time   { return e.CreatedAt }

type fakeDeserializer struct{}

func (f fakeDeserializer) Deserialize(eventType string, data []byte) (domainevent.Event, error) {
	var te testEvent
	if err := json.Unmarshal(data, &te); err != nil {
		return nil, err
	}
	return te, nil
}

func newTestDBClient(t *testing.T) *client.Client {
	t.Helper()

	testCfg, err := config.NewTestDatabaseConfig()
	require.NoError(t, err)

	c, err := client.NewClient(config.DatabaseConfig{
		User:     testCfg.User,
		Password: testCfg.Password,
		Host:     testCfg.Host,
		Port:     testCfg.Port,
		Name:     testCfg.Name,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		err = c.Close()
		require.NoError(t, err)
	})

	return c
}

func beginTxCtx(t *testing.T, dbClient *client.Client) (context.Context, *sqlx.Tx) {
	t.Helper()

	db := dbClient.GetDB()
	tx, err := db.Beginx()
	require.NoError(t, err)

	ctx := transaction.WithTx(context.Background(), tx)

	t.Cleanup(func() { _ = tx.Rollback() })
	return ctx, tx
}

func TestEventStore_SaveEvents(t *testing.T) {
	testAggregateID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := map[string]struct {
		events      []domainevent.Event
		expectError bool
		errorText   string
	}{
		"successful save single event": {
			events: []domainevent.Event{
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TodoListCreated",
					Version:     1,
					CreatedAt:   time.Now(),
				},
			},
			expectError: false,
		},
		"successful save multiple events": {
			events: []domainevent.Event{
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TodoListCreated",
					Version:     1,
					CreatedAt:   time.Now(),
				},
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TodoAdded",
					Version:     2,
					Title:       "Test Todo",
					CreatedAt:   time.Now(),
				},
			},
			expectError: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dbClient := newTestDBClient(t)
			ctx, _ := beginTxCtx(t, dbClient)
			store := eventstore.NewEventStore(fakeDeserializer{})

			err := store.SaveEvents(ctx, testAggregateID, tt.events)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorText != "" {
					require.Contains(t, err.Error(), tt.errorText)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEventStore_SaveEvents_OptimisticLock(t *testing.T) {
	testAggregateID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := map[string]struct {
		events             []domainevent.Event
		expectOptimisticLock bool
	}{
		"no conflict with different versions": {
			events: []domainevent.Event{
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TestEvent",
					Version:     1,
					CreatedAt:   time.Now(),
				},
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TestEvent",
					Version:     2,
					CreatedAt:   time.Now(),
				},
			},
			expectOptimisticLock: false,
		},
		"version conflict same aggregate same version": {
			events: []domainevent.Event{
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TestEvent",
					Version:     1,
					CreatedAt:   time.Now(),
				},
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TestEvent",
					Version:     1,
					CreatedAt:   time.Now(),
				},
			},
			expectOptimisticLock: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dbClient := newTestDBClient(t)
			ctx, _ := beginTxCtx(t, dbClient)
			store := eventstore.NewEventStore(fakeDeserializer{})

			err := store.SaveEvents(ctx, testAggregateID, tt.events)

			if tt.expectOptimisticLock {
				require.Error(t, err)
				var appErr *appErrors.Error
				require.True(t, errors.As(err, &appErr))
				require.Equal(t, appErrors.OptimisticLock, appErr.ErrCode)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEventStore_LoadEvents(t *testing.T) {
	// Fixed aggregate ID for test consistency
	testAggregateID := uuid.MustParse("12345678-1234-1234-1234-123456789012")

	tests := map[string]struct {
		savedEvents   []domainevent.Event
		expectedCount int
		verifyOrder   bool
	}{
		"load events": {
			savedEvents: []domainevent.Event{
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TodoAdded",
					Version:     2,
					Title:       "Second Todo",
					CreatedAt:   time.Now(),
				},
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TodoListCreated",
					Version:     1,
					CreatedAt:   time.Now().Add(-time.Minute),
				},
				testEvent{
					AggregateID: testAggregateID,
					EventID:     uuid.New(),
					Type:        "TodoAdded",
					Version:     3,
					Title:       "Third Todo",
					CreatedAt:   time.Now().Add(time.Minute),
				},
			},
			expectedCount: 3,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			dbClient := newTestDBClient(t)
			ctx, _ := beginTxCtx(t, dbClient)
			store := eventstore.NewEventStore(fakeDeserializer{})
			err := store.SaveEvents(ctx, testAggregateID, tt.savedEvents)
			require.NoError(t, err)

			// Act
			loadedEvents, err := store.LoadEvents(ctx, testAggregateID)

			// Assert
			require.NoError(t, err)
			require.Len(t, loadedEvents, tt.expectedCount)
		})
	}
}

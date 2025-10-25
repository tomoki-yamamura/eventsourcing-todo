package eventstore

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/eventstore/deserializer"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/testutil"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/transaction"
)

func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	return testutil.SetupTestDB(t)
}

func TestEventStoreImpl_SaveEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	tests := map[string]struct {
		events        []event.Event
		expectedCount int
		wantErr       bool
	}{
		"save single event to database": {
			events: []event.Event{
				event.TodoListCreatedEvent{
					AggregateID: uuid.New(),
					UserID:      "user123",
					EventID:     uuid.New(),
					Timestamp:   time.Now(),
					Version:     1,
				},
			},
			expectedCount: 1,
			wantErr:       false,
		},
		"save multiple events to database": {
			events: func() []event.Event {
				aggregateID := uuid.New()
				todoText, _ := value.NewTodoText("Learn Event Sourcing")
				return []event.Event{
					event.TodoListCreatedEvent{
						AggregateID: aggregateID,
						UserID:      "user123",
						EventID:     uuid.New(),
						Timestamp:   time.Now(),
						Version:     1,
					},
					event.TodoAddedEvent{
						AggregateID: aggregateID,
						UserID:      "user123",
						TodoText:    todoText,
						EventID:     uuid.New(),
						Timestamp:   time.Now(),
						Version:     2,
					},
				}
			}(),
			expectedCount: 2,
			wantErr:       false,
		},
		"verify JSON serialization in database": {
			events: []event.Event{
				event.TodoListCreatedEvent{
					AggregateID: uuid.New(),
					UserID:      "test_json_user",
					EventID:     uuid.New(),
					Timestamp:   time.Now(),
					Version:     1,
				},
			},
			expectedCount: 1,
			wantErr:       false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// EventStore作成
			deserializer := deserializer.NewEventDeserializer()
			eventStore := NewEventStore(deserializer)

			// Transaction作成
			txManager := transaction.NewTransaction(db)
			aggregateID := tt.events[0].GetAggregateID()

			// トランザクション内でテスト実行（自動的にrollbackされる）
			err := txManager.RWTx(context.Background(), func(ctx context.Context) error {
				// SaveEventsを実行
				err := eventStore.SaveEvents(ctx, aggregateID, tt.events)
				if (err != nil) != tt.wantErr {
					t.Errorf("SaveEvents() error = %v, wantErr %v", err, tt.wantErr)
				}
				if err != nil {
					return err
				}

				// トランザクション内でデータベースから確認
				tx, err := transaction.GetTx(ctx)
				require.NoError(t, err)

				// DB insertが正しく行われているかを確認
				rows, err := tx.QueryContext(ctx, `
					SELECT event_id, event_type, event_data, version, created_at
					FROM events 
					WHERE aggregate_id = ? 
					ORDER BY version ASC
				`, aggregateID.String())
				require.NoError(t, err)
				defer rows.Close()

				var count int
				for rows.Next() {
					var eventID, eventType string
					var eventData []byte
					var version int
					var createdAt time.Time
					err := rows.Scan(&eventID, &eventType, &eventData, &version, &createdAt)
					require.NoError(t, err)
					
					// DB insertの基本検証
					require.Equal(t, tt.events[count].GetEventType(), eventType)
					require.Equal(t, tt.events[count].GetVersion(), version)
					require.Equal(t, tt.events[count].GetEventID().String(), eventID)
					
					// JSONシリアライゼーションの検証
					require.NotEmpty(t, eventData, "Event data should not be empty")
					require.True(t, json.Valid(eventData), "Event data should be valid JSON")
					
					// created_atが設定されているかの検証
					require.False(t, createdAt.IsZero(), "created_at should be set")
					
					count++
				}
				require.NoError(t, rows.Err())
				require.Equal(t, tt.expectedCount, count, "Expected number of events should be inserted")

				// rollbackさせるためにエラーを返す
				return &rollbackError{}
			})

			// rollbackエラーは期待されるエラー
			_, isRollback := err.(*rollbackError)
			require.True(t, isRollback, "Expected rollback error but got: %v", err)
		})
	}
}

func TestEventStoreImpl_LoadEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	tests := map[string]struct {
		setupEvents   []event.Event
		expectedCount int
		wantErr       bool
	}{
		"load single event": {
			setupEvents: []event.Event{
				event.TodoListCreatedEvent{
					AggregateID: uuid.New(),
					UserID:      "user123",
					EventID:     uuid.New(),
					Timestamp:   time.Now(),
					Version:     1,
				},
			},
			expectedCount: 1,
			wantErr:       false,
		},
		"load multiple events in version order": {
			setupEvents: func() []event.Event {
				aggregateID := uuid.New()
				todoText, _ := value.NewTodoText("Learn Event Sourcing")
				return []event.Event{
					event.TodoListCreatedEvent{
						AggregateID: aggregateID,
						UserID:      "user123",
						EventID:     uuid.New(),
						Timestamp:   time.Now(),
						Version:     1,
					},
					event.TodoAddedEvent{
						AggregateID: aggregateID,
						UserID:      "user123",
						TodoText:    todoText,
						EventID:     uuid.New(),
						Timestamp:   time.Now(),
						Version:     2,
					},
				}
			}(),
			expectedCount: 2,
			wantErr:       false,
		},
		"load no events for non-existent aggregate": {
			setupEvents:   []event.Event{},
			expectedCount: 0,
			wantErr:       false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// EventStore作成
			deserializer := deserializer.NewEventDeserializer()
			eventStore := NewEventStore(deserializer)

			// Transaction作成
			txManager := transaction.NewTransaction(db)

			var aggregateID uuid.UUID
			if len(tt.setupEvents) > 0 {
				aggregateID = tt.setupEvents[0].GetAggregateID()
			} else {
				aggregateID = uuid.New() // 存在しないaggregate ID
			}

			// トランザクション内でテスト実行
			err := txManager.RWTx(context.Background(), func(ctx context.Context) error {
				// セットアップイベントがあれば保存
				if len(tt.setupEvents) > 0 {
					err := eventStore.SaveEvents(ctx, aggregateID, tt.setupEvents)
					require.NoError(t, err)
				}

				// LoadEventsを実行
				loadedEvents, err := eventStore.LoadEvents(ctx, aggregateID)
				if (err != nil) != tt.wantErr {
					t.Errorf("LoadEvents() error = %v, wantErr %v", err, tt.wantErr)
				}
				if err != nil {
					return err
				}

				// 検証
				require.Len(t, loadedEvents, tt.expectedCount)

				for i, evt := range loadedEvents {
					require.Equal(t, tt.setupEvents[i].GetEventType(), evt.GetEventType())
					require.Equal(t, tt.setupEvents[i].GetVersion(), evt.GetVersion())
					require.Equal(t, aggregateID, evt.GetAggregateID())
				}

				// バージョン順序を確認
				for i := 1; i < len(loadedEvents); i++ {
					require.True(t, loadedEvents[i-1].GetVersion() < loadedEvents[i].GetVersion())
				}

				// rollbackさせるためにエラーを返す
				return &rollbackError{}
			})

			// rollbackエラーは期待されるエラー
			_, isRollback := err.(*rollbackError)
			require.True(t, isRollback, "Expected rollback error but got: %v", err)
		})
	}
}

func TestEventStoreImpl_SaveAndLoadEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// EventStore作成
	deserializer := deserializer.NewEventDeserializer()
	eventStore := NewEventStore(deserializer)

	// Transaction作成
	txManager := transaction.NewTransaction(db)

	// テストデータ準備
	aggregateID := uuid.New()
	todoText, err := value.NewTodoText("Learn Event Sourcing with MySQL")
	require.NoError(t, err)

	events := []event.Event{
		event.TodoListCreatedEvent{
			AggregateID: aggregateID,
			UserID:      "test_user",
			EventID:     uuid.New(),
			Timestamp:   time.Now(),
			Version:     1,
		},
		event.TodoAddedEvent{
			AggregateID: aggregateID,
			UserID:      "test_user",
			TodoText:    todoText,
			EventID:     uuid.New(),
			Timestamp:   time.Now(),
			Version:     2,
		},
	}

	// トランザクション内でテスト実行
	err = txManager.RWTx(context.Background(), func(ctx context.Context) error {
		// 保存
		err := eventStore.SaveEvents(ctx, aggregateID, events)
		require.NoError(t, err)

		// ロード
		loadedEvents, err := eventStore.LoadEvents(ctx, aggregateID)
		require.NoError(t, err)

		// 検証
		require.Len(t, loadedEvents, 2)
		require.Equal(t, "TodoListCreatedEvent", loadedEvents[0].GetEventType())
		require.Equal(t, "TodoAddedEvent", loadedEvents[1].GetEventType())
		require.Equal(t, 1, loadedEvents[0].GetVersion())
		require.Equal(t, 2, loadedEvents[1].GetVersion())

		// rollbackさせるためにエラーを返す
		return &rollbackError{}
	})

	// rollbackエラーは期待されるエラー
	_, isRollback := err.(*rollbackError)
	require.True(t, isRollback, "Expected rollback error but got: %v", err)
}

// rollbackError はテスト用のエラー型（トランザクションをrollbackさせるため）
type rollbackError struct{}

func (e *rollbackError) Error() string {
	return "rollback for test cleanup"
}
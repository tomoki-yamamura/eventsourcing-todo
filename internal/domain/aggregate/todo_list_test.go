package aggregate

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
)

func TestTodoListAggregate_CreateTodoList(t *testing.T) {
	tests := map[string]struct {
		userID          string
		expectedVersion int
		expectedEvents  int
		wantErr         bool
	}{
		"valid user ID": {
			userID:          "user123",
			expectedVersion: 1,
			expectedEvents:  1,
			wantErr:         false,
		},
		"empty user ID": {
			userID:          "",
			expectedVersion: 0,
			expectedEvents:  0,
			wantErr:         false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			aggregate := NewTodoListAggregate()
			cmd := command.CreateTodoListCommand{
				UserID: tt.userID,
			}

			// Act
			err := aggregate.ExecuteCreateTodoListCommand(cmd)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.expectedEvents > 0 {
					require.NotEqual(t, uuid.Nil, aggregate.GetAggregateID())
					require.Equal(t, tt.userID, aggregate.GetUserID())
				}
				require.Equal(t, tt.expectedVersion, aggregate.GetVersion())

				uncommittedEvents := aggregate.GetUncommittedEvents()
				require.Len(t, uncommittedEvents, tt.expectedEvents)
				if tt.expectedEvents > 0 {
					require.Equal(t, "TodoListCreatedEvent", uncommittedEvents[0].GetEventType())
				}
			}
		})
	}
}

func TestTodoListAggregate_ExecuteAddTodoCommand(t *testing.T) {
	tests := map[string]struct {
		todoText        string
		expectedVersion int
		expectedEvents  int
		wantErr         bool
	}{
		"valid todo text": {
			todoText:        "Learn Event Sourcing",
			expectedVersion: 2,
			expectedEvents:  2,
			wantErr:         false,
		},
		"long todo text": {
			todoText:        "This is a very long todo item that contains a lot of details about what needs to be done",
			expectedVersion: 2,
			expectedEvents:  2,
			wantErr:         false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			userID := "user123"
			aggregate := NewTodoListAggregate()

			createCmd := command.CreateTodoListCommand{
				UserID: userID,
			}
			err := aggregate.ExecuteCreateTodoListCommand(createCmd)
			require.NoError(t, err)

			todoText, err := value.NewTodoText(tt.todoText)
			require.NoError(t, err)
			addCmd := command.AddTodoCommand{
				AggregateID: aggregate.GetAggregateID(),
				UserID:      userID,
				TodoText:    todoText,
			}

			// Act
			err = aggregate.ExecuteAddTodoCommand(addCmd)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedVersion, aggregate.GetVersion())

				uncommittedEvents := aggregate.GetUncommittedEvents()
				require.Len(t, uncommittedEvents, tt.expectedEvents)
				require.Equal(t, "TodoListCreatedEvent", uncommittedEvents[0].GetEventType())
				require.Equal(t, "TodoAddedEvent", uncommittedEvents[1].GetEventType())
			}
		})
	}
}

func TestTodoListAggregate_ExecuteAddTodoCommand_ExceedsLimit(t *testing.T) {
	tests := map[string]struct {
		todosToAdd         int
		expectedError      string
		expectedEventCount int
		wantErr            bool
	}{
		"add 4th todo exceeds limit": {
			todosToAdd:         4,
			expectedError:      "cannot add more than 3 todos per day",
			expectedEventCount: 4, // Create + 3 Add events
			wantErr:            true,
		},
		"add exactly 3 todos": {
			todosToAdd:         3,
			expectedError:      "",
			expectedEventCount: 4, // Create + 3 Add events
			wantErr:            false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			userID := "user123"
			aggregate := NewTodoListAggregate()

			// First create the todo list
			createCmd := command.CreateTodoListCommand{
				UserID: userID,
			}
			err := aggregate.ExecuteCreateTodoListCommand(createCmd)
			require.NoError(t, err)

			var lastErr error
			// Add todos
			for i := range tt.todosToAdd {
				todoText, err := value.NewTodoText(fmt.Sprintf("Todo %d", i+1))
				require.NoError(t, err)
				cmd := command.AddTodoCommand{
					AggregateID: aggregate.GetAggregateID(),
					UserID:      userID,
					TodoText:    todoText,
				}
				lastErr = aggregate.ExecuteAddTodoCommand(cmd)
				if lastErr != nil {
					break // Stop on first error
				}
			}

			// Assert
			if tt.wantErr {
				require.Error(t, lastErr)
				require.Contains(t, lastErr.Error(), tt.expectedError)
			} else {
				require.NoError(t, lastErr)
			}
			require.Len(t, aggregate.GetUncommittedEvents(), tt.expectedEventCount)
		})
	}
}

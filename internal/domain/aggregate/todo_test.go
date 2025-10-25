package aggregate

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
)

func TestTodoListAggregate_ExecuteAddTodoCommand(t *testing.T) {
	// Arrange
	aggregateID := uuid.New()
	aggregate := NewTodoListAggregate(aggregateID)
	todoText, err := value.NewTodoText("Learn Event Sourcing")
	require.NoError(t, err)
	cmd := command.AddTodoCommand{
		AggregateID: aggregateID,
		TodoText:    todoText,
	}

	// Act
	err = aggregate.ExecuteAddTodoCommand(cmd)
	// Assert
	require.NoError(t, err)
	require.Equal(t, aggregateID, aggregate.GetAggregateID())
	require.Equal(t, 1, aggregate.GetVersion())

	uncommittedEvents := aggregate.GetUncommittedEvents()
	require.Len(t, uncommittedEvents, 1)
	require.Equal(t, "TodoAddedEvent", uncommittedEvents[0].GetEventType())
}

func TestTodoListAggregate_ExecuteAddTodoCommand_EmptyTodo(t *testing.T) {
	// Arrange
	aggregateID := uuid.New()
	aggregate := NewTodoListAggregate(aggregateID)
	_, err := value.NewTodoText("")
	require.Error(t, err)
	require.ErrorIs(t, err, value.ErrTodoTextEmpty)

	// Empty todo should be caught at value object level
	// No need to test aggregate with invalid value object
	require.Empty(t, aggregate.GetUncommittedEvents())
}

func TestTodoListAggregate_ExecuteAddTodoCommand_ExceedsLimit(t *testing.T) {
	// Arrange
	aggregateID := uuid.New()
	aggregate := NewTodoListAggregate(aggregateID)
	
	// Add 3 todos first
	for i := 0; i < 3; i++ {
		todoText, err := value.NewTodoText(fmt.Sprintf("Todo %d", i+1))
		require.NoError(t, err)
		cmd := command.AddTodoCommand{
			AggregateID: aggregateID,
			TodoText:    todoText,
		}
		err = aggregate.ExecuteAddTodoCommand(cmd)
		require.NoError(t, err)
	}

	// Try to add 4th todo
	todoText, err := value.NewTodoText("4th Todo")
	require.NoError(t, err)
	cmd := command.AddTodoCommand{
		AggregateID: aggregateID,
		TodoText:    todoText,
	}

	// Act
	err = aggregate.ExecuteAddTodoCommand(cmd)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot add more than 3 todos per day")
	require.Len(t, aggregate.GetUncommittedEvents(), 3) // Only 3 events should exist
}

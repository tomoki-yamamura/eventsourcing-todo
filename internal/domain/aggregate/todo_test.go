package aggregate

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
)

func TestTodoAggregate_ExecuteAddTodoCommand(t *testing.T) {
	// Arrange
	aggregate := NewTodoAggregate()
	aggregateID := uuid.New()
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

func TestTodoAggregate_ExecuteAddTodoCommand_EmptyTodo(t *testing.T) {
	// Arrange
	aggregate := NewTodoAggregate()
	_, err := value.NewTodoText("")
	require.Error(t, err)
	require.ErrorIs(t, err, value.ErrTodoTextEmpty)

	// Empty todo should be caught at value object level
	// No need to test aggregate with invalid value object
	require.Empty(t, aggregate.GetUncommittedEvents())
}

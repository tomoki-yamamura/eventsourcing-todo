package aggregate

import (
	"testing"

	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/command"
)

func TestTodoAggregate_HandleAddTodoCommand(t *testing.T) {
	// Arrange
	aggregate := NewTodoAggregate()
	aggregateID := uuid.New()
	cmd := command.AddTodoCommand{
		AggregateID: aggregateID,
		Todo:        "Learn Event Sourcing",
	}

	// Act
	err := aggregate.HandleAddTodoCommand(cmd)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if aggregate.GetAggregateID() != aggregateID {
		t.Errorf("Expected aggregate ID %v, got %v", aggregateID, aggregate.GetAggregateID())
	}

	if aggregate.GetVersion() != 1 {
		t.Errorf("Expected version 1, got %d", aggregate.GetVersion())
	}

	uncommittedEvents := aggregate.GetUncommittedEvents()
	if len(uncommittedEvents) != 1 {
		t.Errorf("Expected 1 uncommitted event, got %d", len(uncommittedEvents))
	}

	if uncommittedEvents[0].GetEventType() != "TodoAddedEvent" {
		t.Errorf("Expected TodoAddedEvent, got %s", uncommittedEvents[0].GetEventType())
	}
}

func TestTodoAggregate_HandleAddTodoCommand_EmptyTodo(t *testing.T) {
	// Arrange
	aggregate := NewTodoAggregate()
	cmd := command.AddTodoCommand{
		AggregateID: uuid.New(),
		Todo:        "",
	}

	// Act
	err := aggregate.HandleAddTodoCommand(cmd)

	// Assert
	if err == nil {
		t.Fatal("Expected error for empty todo")
	}

	if len(aggregate.GetUncommittedEvents()) != 0 {
		t.Error("Expected no uncommitted events for invalid command")
	}
}
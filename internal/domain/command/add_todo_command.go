package command

import (
	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
)

type AddTodoCommand struct {
	AggregateID uuid.UUID
	TodoText    value.TodoText
}

package command

import "github.com/google/uuid"

type AddTodoCommand struct {
	AggregateID uuid.UUID
	Todo        string
}

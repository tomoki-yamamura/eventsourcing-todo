package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
)

type TodoItem struct {
	ID        string
	Text      value.TodoText
	CreatedAt time.Time
}

func NewTodoItem(text value.TodoText) *TodoItem {
	return &TodoItem{
		ID:        uuid.New().String(),
		Text:      text,
		CreatedAt: time.Now(),
	}
}

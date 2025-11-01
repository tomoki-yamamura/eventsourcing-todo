package value

import (
	"strings"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/errors"
)

type TodoText string

func NewTodoText(text string) (TodoText, error) {
	trimmed := strings.TrimSpace(text)

	if trimmed == "" {
		return "", errors.NewDomainError(errors.InvalidParameter, "todo_text cannot be empty")
	}

	if len(trimmed) > 256 {
		return "", errors.NewDomainError(errors.InvalidParameter, "todo_text cannot exceed 256 characters")
	}

	return TodoText(trimmed), nil
}

func (t TodoText) String() string {
	return string(t)
}

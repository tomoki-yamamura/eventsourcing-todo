package value

import (
	"strings"
)

type TodoText string

func NewTodoText(text string) (TodoText, error) {
	trimmed := strings.TrimSpace(text)

	if trimmed == "" {
		return "", DomainValidationError{
			Field:   "todo_text",
			Message: "cannot be empty",
		}
	}

	if len(trimmed) > 256 {
		return "", DomainValidationError{
			Field:   "todo_text",
			Message: "cannot exceed 256 characters",
		}
	}

	return TodoText(trimmed), nil
}

func (t TodoText) String() string {
	return string(t)
}

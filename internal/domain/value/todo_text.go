package value

import (
	"errors"
	"strings"
)

var (
	ErrTodoTextEmpty   = errors.New("todo text cannot be empty")
	ErrTodoTextTooLong = errors.New("todo text cannot exceed 256 characters")
)

type TodoText string

func NewTodoText(text string) (TodoText, error) {
	trimmed := strings.TrimSpace(text)

	if trimmed == "" {
		return "", ErrTodoTextEmpty
	}

	if len(trimmed) > 256 {
		return "", ErrTodoTextTooLong
	}

	return TodoText(trimmed), nil
}

func (t TodoText) ToString() string {
	return string(t)
}

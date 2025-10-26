package value

import (
	"errors"
	"strings"
)

var (
	ErrUserIDEmpty   = errors.New("user ID cannot be empty")
	ErrUserIDTooLong = errors.New("user ID cannot exceed 128 characters")
)

type UserID string

func NewUserID(id string) (UserID, error) {
	trimmed := strings.TrimSpace(id)

	if trimmed == "" {
		return "", ErrUserIDEmpty
	}

	if len(trimmed) > 128 {
		return "", ErrUserIDTooLong
	}

	return UserID(trimmed), nil
}

func (u UserID) String() string {
	return string(u)
}

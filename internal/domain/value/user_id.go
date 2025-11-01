package value

import (
	"strings"
)

type UserID string

func NewUserID(id string) (UserID, error) {
	trimmed := strings.TrimSpace(id)

	if trimmed == "" {
		return "", DomainValidationError{
			Field:   "user_id",
			Message: "cannot be empty",
		}
	}

	if len(trimmed) > 128 {
		return "", DomainValidationError{
			Field:   "user_id",
			Message: "cannot exceed 128 characters",
		}
	}

	return UserID(trimmed), nil
}

func (u UserID) String() string {
	return string(u)
}

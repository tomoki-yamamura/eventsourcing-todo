package value

import (
	"strings"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/errors"
)

type UserID string

func NewUserID(id string) (UserID, error) {
	trimmed := strings.TrimSpace(id)

	if trimmed == "" {
		return "", errors.NewDomainError(errors.InvalidParameter, "user_id cannot be empty")
	}

	if len(trimmed) > 128 {
		return "", errors.NewDomainError(errors.InvalidParameter, "user_id cannot exceed 128 characters")
	}

	return UserID(trimmed), nil
}

func (u UserID) String() string {
	return string(u)
}

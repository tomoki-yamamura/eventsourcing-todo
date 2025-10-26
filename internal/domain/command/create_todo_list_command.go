package command

import "github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"

type CreateTodoListCommand struct {
	UserID value.UserID
}

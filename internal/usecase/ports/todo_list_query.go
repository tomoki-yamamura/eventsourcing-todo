package ports

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/dto"
)

type TodoListQuery interface {
	Get(ctx context.Context, aggregateID string) *dto.TodoListViewDTO
}

type TodoListViewRepository interface {
	TodoListQuery
	Save(ctx context.Context, id string, view *dto.TodoListViewDTO) error
}

package readmodelstore

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/readmodelstore/dto"
)

type TodoListReadModelStore interface {
	Get(ctx context.Context, aggregateID string) *dto.TodoListViewDTO
	Save(ctx context.Context, aggregateID string, view *dto.TodoListViewDTO) error
}

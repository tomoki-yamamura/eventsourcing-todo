package todo

import (
	"context"
	"sync"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/errors"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/readmodelstore"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/readmodelstore/dto"
)

type InMemoryTodoListViewRepository struct {
	mu   sync.RWMutex
	data map[string]*dto.TodoListViewDTO
}

func NewInMemoryTodoListViewRepository() readmodelstore.TodoListReadModelStore {
	return &InMemoryTodoListViewRepository{
		data: make(map[string]*dto.TodoListViewDTO),
	}
}

func (r *InMemoryTodoListViewRepository) Get(ctx context.Context, aggregateID string) (*dto.TodoListViewDTO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	view := r.data[aggregateID]
	if view == nil {
		return nil, errors.NotFound.New("todo list not found")
	}

	return r.cloneView(view), nil
}

func (r *InMemoryTodoListViewRepository) Save(ctx context.Context, aggregateID string, view *dto.TodoListViewDTO) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[aggregateID] = r.cloneView(view)
	return nil
}

func (r *InMemoryTodoListViewRepository) cloneView(view *dto.TodoListViewDTO) *dto.TodoListViewDTO {
	if view == nil {
		return nil
	}

	items := make([]dto.TodoItemViewDTO, len(view.Items))
	copy(items, view.Items)

	return &dto.TodoListViewDTO{
		AggregateID: view.AggregateID,
		UserID:      view.UserID,
		Items:       items,
		Version:     view.Version,
		UpdatedAt:   view.UpdatedAt,
	}
}

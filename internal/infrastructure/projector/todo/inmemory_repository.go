package todo

import (
	"sync"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/dto"
)

type InMemoryTodoListViewRepository struct {
	mu   sync.RWMutex
	data map[string]*dto.TodoListViewDTO
}

func NewInMemoryTodoListViewRepository() ports.TodoListViewRepository {
	return &InMemoryTodoListViewRepository{
		data: make(map[string]*dto.TodoListViewDTO),
	}
}

func (r *InMemoryTodoListViewRepository) Get(aggregateID string) *dto.TodoListViewDTO {
	r.mu.RLock()
	defer r.mu.RUnlock()

	view := r.data[aggregateID]
	if view == nil {
		return nil
	}

	return r.cloneView(view)
}

func (r *InMemoryTodoListViewRepository) Save(aggregateID string, view *dto.TodoListViewDTO) error {
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

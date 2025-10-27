package todo

import (
	"sync"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports"
)

type InMemoryTodoListViewRepository struct {
	mu   sync.RWMutex
	data map[string]*TodoListViewDTO
}

func NewInMemoryTodoListViewRepository() ports.Query[*TodoListViewDTO] {
	return &InMemoryTodoListViewRepository{
		data: make(map[string]*TodoListViewDTO),
	}
}

func (r *InMemoryTodoListViewRepository) Get(id string) *TodoListViewDTO {
	r.mu.RLock()
	defer r.mu.RUnlock()

	view := r.data[id]
	if view == nil {
		return nil
	}

	return r.cloneView(view)
}

func (r *InMemoryTodoListViewRepository) Save(aggregateID string, view *TodoListViewDTO) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[aggregateID] = r.cloneView(view)
	return nil
}

func (r *InMemoryTodoListViewRepository) Delete(aggregateID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, aggregateID)
	return nil
}

func (r *InMemoryTodoListViewRepository) cloneView(view *TodoListViewDTO) *TodoListViewDTO {
	if view == nil {
		return nil
	}

	items := make([]TodoItemViewDTO, len(view.Items))
	copy(items, view.Items)

	return &TodoListViewDTO{
		AggregateID: view.AggregateID,
		UserID:      view.UserID,
		Items:       items,
		Version:     view.Version,
		UpdatedAt:   view.UpdatedAt,
	}
}

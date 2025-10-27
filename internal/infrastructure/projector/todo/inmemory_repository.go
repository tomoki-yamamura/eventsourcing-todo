package todo

import "sync"

type InMemoryTodoListViewRepository struct {
	mu   sync.RWMutex
	data map[string]*TodoListView
}

func NewInMemoryTodoListViewRepository() TodoListViewRepository {
	return &InMemoryTodoListViewRepository{
		data: make(map[string]*TodoListView),
	}
}

func (r *InMemoryTodoListViewRepository) Get(aggregateID string) *TodoListView {
	r.mu.RLock()
	defer r.mu.RUnlock()

	view := r.data[aggregateID]
	if view == nil {
		return nil
	}

	return r.cloneView(view)
}

func (r *InMemoryTodoListViewRepository) Save(aggregateID string, view *TodoListView) error {
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

func (r *InMemoryTodoListViewRepository) cloneView(view *TodoListView) *TodoListView {
	if view == nil {
		return nil
	}

	items := make([]TodoItemView, len(view.Items))
	copy(items, view.Items)

	return &TodoListView{
		AggregateID: view.AggregateID,
		UserID:      view.UserID,
		Items:       items,
		Version:     view.Version,
		UpdatedAt:   view.UpdatedAt,
	}
}

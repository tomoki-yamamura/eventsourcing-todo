package todo

import "time"

type TodoListView struct {
	AggregateID string
	UserID      string
	Items       []TodoItemView
	Version     int
	UpdatedAt   time.Time
}

type TodoItemView struct {
	Text string
}

type TodoListViewRepository interface {
	Get(aggregateID string) *TodoListView
	Save(aggregateID string, view *TodoListView) error
	Delete(aggregateID string) error
}

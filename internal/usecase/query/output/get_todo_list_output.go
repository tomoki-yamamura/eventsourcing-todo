package output

import "time"

type GetTodoListOutput struct {
	AggregateID string
	UserID      string
	Items       []TodoItem
	UpdatedAt   time.Time
}

type TodoItem struct {
	Text string
}

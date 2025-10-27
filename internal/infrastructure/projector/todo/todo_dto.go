package todo

import "time"

// TodoListViewDTO - 読み取りモデル用DTO
type TodoListViewDTO struct {
	AggregateID string
	UserID      string
	Items       []TodoItemViewDTO
	Version     int
	UpdatedAt   time.Time
}

// TodoItemViewDTO - TodoItemの読み取りモデル用DTO
type TodoItemViewDTO struct {
	Text string
}

package dto

import "time"

type TodoListViewDTO struct {
	AggregateID string
	UserID      string
	Items       []TodoItemViewDTO
	Version     int
	UpdatedAt   time.Time
}

type TodoItemViewDTO struct {
	Text string
}

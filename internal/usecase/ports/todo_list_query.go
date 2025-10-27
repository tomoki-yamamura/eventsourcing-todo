package ports

import "github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/dto"

type TodoListQuery interface {
	Get(aggregateID string) *dto.TodoListViewDTO
}

type TodoListViewRepository interface {
	TodoListQuery
	Save(id string, view *dto.TodoListViewDTO) error
}

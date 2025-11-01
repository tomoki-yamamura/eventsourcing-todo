package presenter

import (
	"context"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/presenter/viewmodel"
)

type TodoListView interface {
	Render(ctx context.Context, vm *viewmodel.TodoListVM, status int, err error)
}

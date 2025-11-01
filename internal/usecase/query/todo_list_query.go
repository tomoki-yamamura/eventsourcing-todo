package query

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/presenter"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/readmodelstore"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/input"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/output"
)

type TodoListQueryInterface interface {
	Execute(ctx context.Context, input *input.GetTodoListInput, out presenter.TodoListPresenter) error
}

type TodoListQuery struct {
	store readmodelstore.TodoListReadModelStore
}

func NewTodoListQuery(store readmodelstore.TodoListReadModelStore) TodoListQueryInterface {
	return &TodoListQuery{
		store: store,
	}
}

func (u *TodoListQuery) Execute(ctx context.Context, input *input.GetTodoListInput, out presenter.TodoListPresenter) error {
	view := u.store.Get(ctx, input.AggregateID)
	if view == nil {
		notFoundErr := value.NotFoundError{
			Resource: "todo list",
			ID:       input.AggregateID,
		}
		return out.PresentNotFound(ctx, notFoundErr)
	}

	items := make([]output.TodoItem, 0, len(view.Items))
	for _, item := range view.Items {
		items = append(items, output.TodoItem{
			Text: item.Text,
		})
	}

	outputData := &output.GetTodoListOutput{
		AggregateID: view.AggregateID,
		UserID:      view.UserID,
		Items:       items,
		UpdatedAt:   view.UpdatedAt,
	}

	return out.Present(ctx, outputData)
}

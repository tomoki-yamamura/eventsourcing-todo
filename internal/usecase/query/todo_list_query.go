package query

import (
	"context"
	"fmt"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/input"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/output"
)

type TodoListQueryInterface interface {
	Query(ctx context.Context, input *input.GetTodoListInput) (*output.GetTodoListOutput, error)
}

type TodoListQuery struct {
	viewQuery ports.TodoListQuery
}

func NewTodoListQuery(viewQuery ports.TodoListQuery) TodoListQueryInterface {
	return &TodoListQuery{
		viewQuery: viewQuery,
	}
}

func (u *TodoListQuery) Query(ctx context.Context, input *input.GetTodoListInput) (*output.GetTodoListOutput, error) {
	view := u.viewQuery.Get(ctx, input.AggregateID)
	if view == nil {
		return nil, fmt.Errorf("todo list not found")
	}

	items := make([]output.TodoItem, 0, len(view.Items))
	for _, item := range view.Items {
		items = append(items, output.TodoItem{
			Text: item.Text,
		})
	}

	return &output.GetTodoListOutput{
		AggregateID: view.AggregateID,
		UserID:      view.UserID,
		Items:       items,
	}, nil
}

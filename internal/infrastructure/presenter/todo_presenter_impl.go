package presenter

import (
	"context"
	"net/http"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/presenter/viewmodel"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/view"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/presenter"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/output"
)

type HTTPTodoListPresenter struct {
	view view.TodoListView
}

func NewHTTPTodoListPresenter(view view.TodoListView) presenter.TodoListPresenter {
	return &HTTPTodoListPresenter{view: view}
}

func (p *HTTPTodoListPresenter) Present(ctx context.Context, out *output.GetTodoListOutput) error {
	var items []viewmodel.TodoItem
	for _, it := range out.Items {
		items = append(items, viewmodel.TodoItem{Text: it.Text})
	}
	vm := &viewmodel.TodoListVM{
		AggregateID: out.AggregateID,
		UserID:      out.UserID,
		Items:       items,
	}
	p.view.Render(ctx, vm, http.StatusOK, nil)
	return nil
}

func (p *HTTPTodoListPresenter) PresentNotFound(ctx context.Context, aggregateID string) error {
	p.view.Render(ctx, nil, http.StatusNotFound, nil)
	return nil
}

func (p *HTTPTodoListPresenter) PresentError(ctx context.Context, err error) error {
	p.view.Render(ctx, nil, http.StatusInternalServerError, err)
	return nil
}

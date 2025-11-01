package view

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/presenter/viewmodel"
)

type TodoListView interface {
	Render(ctx context.Context, vm *viewmodel.TodoListVM, status int, err error)
}

type HTTPTodoListView struct {
	w http.ResponseWriter
}

func NewHTTPTodoListView(w http.ResponseWriter) TodoListView {
	return &HTTPTodoListView{w: w}
}

func (v *HTTPTodoListView) Render(ctx context.Context, vm *viewmodel.TodoListVM, status int, err error) {
	v.w.Header().Set("Content-Type", "application/json")
	v.w.WriteHeader(status)

	if err != nil {
		_ = json.NewEncoder(v.w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if vm != nil {
		_ = json.NewEncoder(v.w).Encode(vm)
	}
}

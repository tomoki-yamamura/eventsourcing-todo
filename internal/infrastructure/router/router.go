package router

import (
	"net/http"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler"
)

type Router struct {
	todoHandler *handler.TodoHandler
}

func NewRouter(todoHandler *handler.TodoHandler) *Router {
	return &Router{
		todoHandler: todoHandler,
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// POST /todo-lists - Create new todo list
	mux.HandleFunc("/todo-lists", r.todoHandler.CreateTodoList)

	// POST /todo-lists/{aggregate_id}/todos - Add todo to existing list
	mux.HandleFunc("/todo-lists/", r.todoHandler.AddTodo)

	return mux
}
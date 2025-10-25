package router

import (
	"github.com/gorilla/mux"
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

func (r *Router) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// POST /todo-lists - Create new todo list
	router.HandleFunc("/todo-lists", r.todoHandler.CreateTodoList).Methods("POST")

	// GET /todo-lists/{aggregate_id}/todos - Get todos in the list
	router.HandleFunc("/todo-lists/{aggregate_id}/todos", r.todoHandler.GetTodoList).Methods("GET")

	// POST /todo-lists/{aggregate_id}/todos - Add todo to existing list
	router.HandleFunc("/todo-lists/{aggregate_id}/todos", r.todoHandler.AddTodo).Methods("POST")

	return router
}
package router

import (
	"github.com/gorilla/mux"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/query"
)

type Router struct {
	commandHandler *command.TodoCommandHandler
	queryHandler   *query.TodoQueryHandler
}

func NewRouter(commandHandler *command.TodoCommandHandler, queryHandler *query.TodoQueryHandler) *Router {
	return &Router{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
	}
}

func (r *Router) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/todo-lists", r.commandHandler.CreateTodoList).Methods("POST")
	router.HandleFunc("/todo-lists/{aggregate_id}/todos", r.commandHandler.AddTodo).Methods("POST")

	router.HandleFunc("/todo-lists/{aggregate_id}/todos", r.queryHandler.GetTodoList).Methods("GET")

	return router
}

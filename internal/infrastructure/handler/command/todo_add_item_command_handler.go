package command

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/request"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command/input"
)

type TodoAddItemCommandHandler struct {
	addCommand command.TodoAddItemCommandInterface
}

func NewTodoAddItemCommandHandler(addCommand command.TodoAddItemCommandInterface) *TodoAddItemCommandHandler {
	return &TodoAddItemCommandHandler{
		addCommand: addCommand,
	}
}

func (h *TodoAddItemCommandHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aggregateID := vars["aggregate_id"]

	var req request.AddTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	usecaseInput := &input.AddTodoInput{
		AggregateID: aggregateID,
		UserID:      req.UserID,
		Todo:        req.Text,
	}

	err := h.addCommand.Execute(r.Context(), usecaseInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Todo item added successfully"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

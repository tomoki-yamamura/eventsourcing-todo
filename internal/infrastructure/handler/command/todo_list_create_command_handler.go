package command

import (
	"encoding/json"
	"net/http"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/request"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/response"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command/input"
)

type TodoListCreateCommandHandler struct {
	createCommand command.TodoListCreateCommandInterface
}

func NewTodoListCreateCommandHandler(createCommand command.TodoListCreateCommandInterface) *TodoListCreateCommandHandler {
	return &TodoListCreateCommandHandler{
		createCommand: createCommand,
	}
}

func (h *TodoListCreateCommandHandler) CreateTodoList(w http.ResponseWriter, r *http.Request) {
	var req request.CreateTodoListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	usecaseInput := &input.CreateTodoListInput{
		UserID: req.UserID,
	}

	err := h.createCommand.Execute(r.Context(), usecaseInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := response.CreateTodoListResponse{
		Message: "Todo list created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

package command

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/request"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/response"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command/input"
)

type TodoCommandHandler struct {
	commandUseCase command.TodoCommandUseCaseInterface
}

func NewTodoCommandHandler(commandUseCase command.TodoCommandUseCaseInterface) *TodoCommandHandler {
	return &TodoCommandHandler{
		commandUseCase: commandUseCase,
	}
}

func (h *TodoCommandHandler) CreateTodoList(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.commandUseCase.CreateTodoList(r.Context(), usecaseInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := response.CreateTodoListResponse{
		AggregateID: result.AggregateID,
		Message:     "Todo list created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *TodoCommandHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.commandUseCase.AddTodo(r.Context(), usecaseInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

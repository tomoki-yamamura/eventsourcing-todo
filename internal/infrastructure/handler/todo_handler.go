package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/request"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/handler/response"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/input"
)

type TodoHandler struct {
	todoUseCase usecase.TodoUseCaseInterface
}

func NewTodoHandler(todoUseCase usecase.TodoUseCaseInterface) *TodoHandler {
	return &TodoHandler{
		todoUseCase: todoUseCase,
	}
}

func (h *TodoHandler) CreateTodoList(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.todoUseCase.CreateTodoList(r.Context(), usecaseInput)
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

func (h *TodoHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.todoUseCase.AddTodo(r.Context(), usecaseInput)
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

func (h *TodoHandler) GetTodoList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aggregateID := vars["aggregate_id"]

	if aggregateID == "" {
		http.Error(w, "aggregate_id is required", http.StatusBadRequest)
		return
	}

	input := &input.GetTodoListInput{
		AggregateID: aggregateID,
	}

	result, err := h.todoUseCase.GetTodoList(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

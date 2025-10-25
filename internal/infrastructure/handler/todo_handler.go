package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase"
)

type TodoHandler struct {
	todoUseCase usecase.TodoUseCaseInterface
}

func NewTodoHandler(todoUseCase usecase.TodoUseCaseInterface) *TodoHandler {
	return &TodoHandler{
		todoUseCase: todoUseCase,
	}
}

type CreateTodoListRequest struct {
	UserID string `json:"user_id"`
}

type CreateTodoListResponse struct {
	AggregateID string `json:"aggregate_id"`
	Message     string `json:"message"`
}

type AddTodoRequest struct {
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

type AddTodoResponse struct {
	Message string `json:"message"`
}

func (h *TodoHandler) CreateTodoList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTodoListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	aggregateID, err := h.todoUseCase.CreateTodoList(r.Context(), req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := CreateTodoListResponse{
		AggregateID: aggregateID,
		Message:     "Todo list created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *TodoHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract aggregate_id from URL path: /todo-lists/{aggregate_id}/todos
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "todo-lists" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	aggregateID := pathParts[1]

	var req AddTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	if err := h.todoUseCase.AddTodo(r.Context(), aggregateID, req.UserID, req.Text); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := AddTodoResponse{
		Message: "Todo added successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
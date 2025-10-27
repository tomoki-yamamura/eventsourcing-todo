package query

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/input"
)

type TodoQueryHandler struct {
	queryUseCase query.TodoQueryUseCaseInterface
}

func NewTodoQueryHandler(queryUseCase query.TodoQueryUseCaseInterface) *TodoQueryHandler {
	return &TodoQueryHandler{
		queryUseCase: queryUseCase,
	}
}

func (h *TodoQueryHandler) GetTodoList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aggregateID := vars["aggregate_id"]

	if aggregateID == "" {
		http.Error(w, "aggregate_id is required", http.StatusBadRequest)
		return
	}

	usecaseInput := &input.GetTodoListInput{
		AggregateID: aggregateID,
	}

	result, err := h.queryUseCase.GetTodoList(r.Context(), usecaseInput)
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

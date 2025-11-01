package view

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/presenter/viewmodel"
)

type HTTPCommandResultView struct {
	writer http.ResponseWriter
}

func NewHTTPCommandResultView(w http.ResponseWriter) *HTTPCommandResultView {
	return &HTTPCommandResultView{
		writer: w,
	}
}

func (v *HTTPCommandResultView) Render(ctx context.Context, vm *viewmodel.CommandResultViewModel, status int, err error) error {
	v.writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		v.writer.WriteHeader(status)
		return json.NewEncoder(v.writer).Encode(map[string]interface{}{
			"status":     "error",
			"message":    err.Error(),
			"executedAt": vm.ExecutedAt,
		})
	}

	v.writer.WriteHeader(status)
	return json.NewEncoder(v.writer).Encode(vm)
}

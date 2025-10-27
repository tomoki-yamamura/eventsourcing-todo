package projector

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/bus"
)

type Projector interface {
	Start(ctx context.Context, eventBus bus.EventBus) error
}

type ViewRepository[T any] interface {
	Get(id string) T
	Save(id string, view T) error
	Delete(id string) error
}

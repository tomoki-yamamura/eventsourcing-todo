package ports

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
)

type Projector interface {
	Handle(ctx context.Context, e event.Event) error
	Start(ctx context.Context, bus EventSubscriber) error
}

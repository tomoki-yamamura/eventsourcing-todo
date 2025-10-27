package ports

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
)

type EventPublisher interface {
	Publish(ctx context.Context, events ...event.Event) error
}

type EventSubscriber interface {
	Subscribe(handler func(context.Context, event.Event) error) error
}

type EventBus interface {
	EventPublisher
	EventSubscriber
}

type Projector interface {
	EventSubscriber
	Start(ctx context.Context, eventBus EventSubscriber) error
}

package bus

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports"
)

type InMemoryEventBus struct {
	handlers []func(context.Context, event.Event) error
}

func NewInMemoryEventBus() ports.EventBus {
	return &InMemoryEventBus{
		handlers: make([]func(context.Context, event.Event) error, 0),
	}
}

func (b *InMemoryEventBus) Publish(ctx context.Context, events ...event.Event) error {
	for _, evt := range events {
		for _, handler := range b.handlers {
			if err := handler(ctx, evt); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *InMemoryEventBus) Subscribe(handler func(context.Context, event.Event) error) error {
	b.handlers = append(b.handlers, handler)
	return nil
}

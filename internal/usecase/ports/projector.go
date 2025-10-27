package ports

import "context"

type Projector interface {
	EventSubscriber
	Start(ctx context.Context, eventBus EventSubscriber) error
}

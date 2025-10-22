package container

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/config"
)

type Container struct {
	Cfg *config.Config
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Inject(ctx context.Context, cfg *config.Config) error {
	return nil
}

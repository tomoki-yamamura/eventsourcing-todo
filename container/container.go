package container

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/config"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/client"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/eventstore"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/transaction"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase"
)

type Container struct {
	// Config
	Cfg *config.Config

	DatabaseClient repository.DatabaseClient

	// Repository layer
	Transaction repository.Transaction
	EventStore  repository.EventStore

	// Use case layer
	TodoUseCase usecase.TodoUseCaseInterface
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Inject(ctx context.Context, cfg *config.Config) error {
	c.Cfg = cfg

	databaseClient, err := client.NewClient(cfg.DatabaseConfig)
	if err != nil {
		return err
	}
	c.DatabaseClient = databaseClient

	// Repository layer
	c.Transaction = transaction.NewTransaction(c.DatabaseClient.GetDB())
	c.EventStore = eventstore.NewEventStore()

	// Use case layer
	c.TodoUseCase = usecase.NewTodoUseCase(c.Transaction, c.EventStore)

	return nil
}

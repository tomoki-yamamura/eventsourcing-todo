package container

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/config"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/bus"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/client"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/eventstore"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/eventstore/deserializer"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/database/transaction"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/projector"
	commandUseCase "github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command"
	queryUseCase "github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query"
)

type Container struct {
	// Config
	Cfg *config.Config

	DatabaseClient repository.DatabaseClient

	// Repository layer
	Transaction  repository.Transaction
	EventStore   repository.EventStore
	Deserializer repository.EventDeserializer

	// Event Bus and Projector
	EventBus  bus.EventBus
	Projector *projector.InMemTodoProjector

	// Use case layer (CQRS)
	CommandUseCase commandUseCase.TodoCommandUseCaseInterface
	QueryUseCase   queryUseCase.TodoQueryUseCaseInterface
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
	c.Deserializer = deserializer.NewEventDeserializer()
	c.EventStore = eventstore.NewEventStore(c.Deserializer)

	// Event Bus and Projector
	c.EventBus = bus.NewInMemoryEventBus()
	c.Projector = projector.NewInMemTodoProjector()

	// Start projector (subscribe to event bus)
	if err := c.Projector.Start(ctx, c.EventBus); err != nil {
		return err
	}

	// Use case layer (CQRS)
	c.CommandUseCase = commandUseCase.NewTodoCommandUseCase(c.Transaction, c.EventStore, c.EventBus)
	c.QueryUseCase = queryUseCase.NewTodoQueryUseCase(c.Transaction, c.EventStore, c.Projector)

	return nil
}

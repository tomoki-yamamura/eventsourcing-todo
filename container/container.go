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
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/projector/todo"
	commandUseCase "github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/gateway"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/readmodelstore"
	queryUseCase "github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query"
)

type Container struct {
	// Config
	Cfg *config.Config

	// Repository layer
	Transaction  repository.Transaction
	EventStore   repository.EventStore
	Deserializer repository.EventDeserializer

	// Gateway implementation
	EventBus      gateway.EventBus
	TodoProjector gateway.Projector
	TodoViewRepo  readmodelstore.TodoListStore

	// Use case layer (CQRS)
	TodoListCreateCommand commandUseCase.TodoListCreateCommandInterface
	TodoAddItemCommand    commandUseCase.TodoAddItemCommandInterface
	QueryUseCase          queryUseCase.TodoListQueryInterface
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

	// Repository layer
	c.Transaction = transaction.NewTransaction(databaseClient.GetDB())
	c.Deserializer = deserializer.NewEventDeserializer()
	c.EventStore = eventstore.NewEventStore(c.Deserializer)

	// Event Bus and Projector
	c.EventBus = bus.NewInMemoryEventBus()
	viewRepo := todo.NewInMemoryTodoListViewRepository()
	c.TodoViewRepo = viewRepo
	c.TodoProjector = todo.NewTodoProjector(viewRepo)

	// Use case layer (CQRS)
	c.TodoListCreateCommand = commandUseCase.NewTodoListCreateCommand(c.Transaction, c.EventStore, c.EventBus)
	c.TodoAddItemCommand = commandUseCase.NewTodoAddItemCommand(c.Transaction, c.EventStore, c.EventBus)
	c.QueryUseCase = queryUseCase.NewTodoListQuery(c.TodoViewRepo)

	return nil
}

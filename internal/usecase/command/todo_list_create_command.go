package command

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/aggregate"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command/input"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/gateway"
)

type TodoListCreateCommandInterface interface {
	Execute(ctx context.Context, input *input.CreateTodoListInput) error
}

type TodoListCreateCommand struct {
	tx         repository.Transaction
	eventStore repository.EventStore
	eventBus   gateway.EventPublisher
}

func NewTodoListCreateCommand(tx repository.Transaction, eventStore repository.EventStore, eventBus gateway.EventPublisher) TodoListCreateCommandInterface {
	return &TodoListCreateCommand{
		tx:         tx,
		eventStore: eventStore,
		eventBus:   eventBus,
	}
}

func (u *TodoListCreateCommand) Execute(ctx context.Context, input *input.CreateTodoListInput) error {
	return u.tx.RWTx(ctx, func(ctx context.Context) error {
		userID, err := value.NewUserID(input.UserID)
		if err != nil {
			return err
		}

		cmd := command.CreateTodoListCommand{
			UserID: userID,
		}

		todoList := aggregate.NewTodoListAggregate()
		if err := todoList.ExecuteCreateTodoListCommand(cmd); err != nil {
			return err
		}

		if err := u.eventStore.SaveEvents(ctx, todoList.GetAggregateID(), todoList.GetUncommittedEvents()); err != nil {
			return err
		}

		evs := todoList.GetUncommittedEvents()
		u.tx.AfterCommit(func() error {
			return u.eventBus.Publish(context.Background(), evs...)
		})

		todoList.MarkEventsAsCommitted()

		return nil
	})
}

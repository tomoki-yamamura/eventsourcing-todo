package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/aggregate"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command/input"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports"
)

type TodoAddItemCommandInterface interface {
	Execute(ctx context.Context, input *input.AddTodoInput) error
}

type TodoAddItemCommand struct {
	tx         repository.Transaction
	eventStore repository.EventStore
	eventBus   ports.EventPublisher
}

func NewTodoAddItemCommand(tx repository.Transaction, eventStore repository.EventStore, eventBus ports.EventPublisher) TodoAddItemCommandInterface {
	return &TodoAddItemCommand{
		tx:         tx,
		eventStore: eventStore,
		eventBus:   eventBus,
	}
}

func (u *TodoAddItemCommand) Execute(ctx context.Context, input *input.AddTodoInput) error {
	maxRetries := 3
	var lastErr error

	for attempt := range maxRetries {
		err := u.tx.RWTx(ctx, func(ctx context.Context) error {
			todoText, err := value.NewTodoText(input.Todo)
			if err != nil {
				return err
			}

			aggregateUUID, err := uuid.Parse(input.AggregateID)
			if err != nil {
				return err
			}

			events, err := u.eventStore.LoadEvents(ctx, aggregateUUID)
			if err != nil {
				return err
			}

			if len(events) == 0 {
				return fmt.Errorf("todo list not found")
			}

			todoList := aggregate.NewTodoListAggregate()
			if err := todoList.Hydration(events); err != nil {
				return err
			}

			userIDVO, err := value.NewUserID(input.UserID)
			if err != nil {
				return err
			}

			cmd := command.AddTodoCommand{
				AggregateID: aggregateUUID,
				UserID:      userIDVO,
				TodoText:    todoText,
			}

			if err := todoList.ExecuteAddTodoCommand(cmd); err != nil {
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

		if err == nil {
			return nil
		}

		lastErr = err

		if isOptimisticLockError(err) && attempt < maxRetries-1 {
			waitTime := time.Duration(attempt+1) * 10 * time.Millisecond
			time.Sleep(waitTime)
			continue
		}

		return err
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

func isOptimisticLockError(err error) bool {
	return strings.Contains(err.Error(), "optimistic lock error")
}
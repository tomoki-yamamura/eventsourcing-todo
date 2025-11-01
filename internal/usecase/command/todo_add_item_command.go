package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/aggregate"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/value"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/command/input"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/gateway"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/presenter"
)

type TodoAddItemCommandInterface interface {
	Execute(ctx context.Context, input *input.AddTodoInput, out presenter.CommandResultPresenter) error
}

type TodoAddItemCommand struct {
	tx         repository.Transaction
	eventStore repository.EventStore
	eventBus   gateway.EventPublisher
}

func NewTodoAddItemCommand(tx repository.Transaction, eventStore repository.EventStore, eventBus gateway.EventPublisher) TodoAddItemCommandInterface {
	return &TodoAddItemCommand{
		tx:         tx,
		eventStore: eventStore,
		eventBus:   eventBus,
	}
}

func (u *TodoAddItemCommand) Execute(ctx context.Context, input *input.AddTodoInput, out presenter.CommandResultPresenter) error {
	maxRetries := 3
	var lastErr error

	for attempt := range maxRetries {
		var aggregateID string
		var version int
		var events []event.Event
		
		err := u.tx.RWTx(ctx, func(ctx context.Context) error {
			todoText, err := value.NewTodoText(input.Todo)
			if err != nil {
				return err
			}

			aggregateUUID, err := uuid.Parse(input.AggregateID)
			if err != nil {
				return err
			}

			loadedEvents, err := u.eventStore.LoadEvents(ctx, aggregateUUID)
			if err != nil {
				return err
			}

			if len(loadedEvents) == 0 {
				return value.NotFoundError{
					Resource: "todo list",
					ID:       input.AggregateID,
				}
			}

			todoList := aggregate.NewTodoListAggregate()
			if err := todoList.Hydration(loadedEvents); err != nil {
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

			aggregateID = todoList.GetAggregateID().String()
			version = todoList.GetVersion()
			events = todoList.GetUncommittedEvents()

			evs := todoList.GetUncommittedEvents()
			u.tx.AfterCommit(func() error {
				return u.eventBus.Publish(context.Background(), evs...)
			})

			todoList.MarkEventsAsCommitted()

			return nil
		})

		if err == nil {
			return out.PresentSuccess(ctx, aggregateID, version, events)
		}

		lastErr = err

		if isOptimisticLockError(err) && attempt < maxRetries-1 {
			waitTime := time.Duration(attempt+1) * 10 * time.Millisecond
			time.Sleep(waitTime)
			continue
		}

		return out.PresentError(ctx, err)
	}

	return out.PresentError(ctx, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr))
}

func isOptimisticLockError(err error) bool {
	return strings.Contains(err.Error(), "optimistic lock error")
}

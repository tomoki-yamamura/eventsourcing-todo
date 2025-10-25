package usecase

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
)

type TodoUseCaseInterface interface {
	CreateTodoList(ctx context.Context, userID string) (string, error)
	AddTodo(ctx context.Context, aggregateID string, userID string, todo string) error
}

type TodoUseCase struct {
	tx         repository.Transaction
	eventStore repository.EventStore
}

func NewTodoUseCase(tx repository.Transaction, eventStore repository.EventStore) TodoUseCaseInterface {
	return &TodoUseCase{
		tx:         tx,
		eventStore: eventStore,
	}
}

func (u *TodoUseCase) CreateTodoList(ctx context.Context, userID string) (string, error) {
	var aggregateID string
	err := u.tx.RWTx(ctx, func(ctx context.Context) error {
		cmd := command.CreateTodoListCommand{
			UserID: userID,
		}

		todoList := aggregate.NewTodoListAggregate()
		if err := todoList.ExecuteCreateTodoListCommand(cmd); err != nil {
			return fmt.Errorf("failed to create todo list: %w", err)
		}

		aggregateID = todoList.GetAggregateID().String()

		if err := u.eventStore.SaveEvents(ctx, todoList.GetAggregateID(), todoList.GetUncommittedEvents()); err != nil {
			return fmt.Errorf("failed to save events: %w", err)
		}
		
		todoList.MarkEventsAsCommitted()

		return nil
	})

	if err != nil {
		return "", err
	}

	return aggregateID, nil
}

func (u *TodoUseCase) AddTodo(ctx context.Context, aggregateID string, userID string, todo string) error {
	return u.addTodoWithRetry(ctx, aggregateID, userID, todo, 3)
}

func (u *TodoUseCase) addTodoWithRetry(ctx context.Context, aggregateID string, userID string, todo string, maxRetries int) error {
	var lastErr error
	
	for attempt := range maxRetries {
		err := u.tx.RWTx(ctx, func(ctx context.Context) error {
			todoText, err := value.NewTodoText(todo)
			if err != nil {
				return err
			}

			aggregateUUID, err := uuid.Parse(aggregateID)
			if err != nil {
				return fmt.Errorf("invalid aggregate ID: %w", err)
			}

			events, err := u.eventStore.LoadEvents(ctx, aggregateUUID)
			if err != nil {
				return fmt.Errorf("failed to load events: %w", err)
			}
			
			if len(events) == 0 {
				return fmt.Errorf("todo list not found")
			}
			
			todoList := aggregate.NewTodoListAggregate()
			if err := todoList.Hydration(events); err != nil {
				return fmt.Errorf("failed to hydrate aggregate: %w", err)
			}
			
			cmd := command.AddTodoCommand{
				AggregateID: aggregateUUID,
				UserID:      userID,
				TodoText:    todoText,
			}
			
			if err := todoList.ExecuteAddTodoCommand(cmd); err != nil {
				return fmt.Errorf("failed to handle add todo command: %w", err)
			}

			if err := u.eventStore.SaveEvents(ctx, todoList.GetAggregateID(), todoList.GetUncommittedEvents()); err != nil {
				return fmt.Errorf("failed to save events: %w", err)
			}
			
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

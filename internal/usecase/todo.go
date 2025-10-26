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
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/input"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/output"
)

type TodoUseCaseInterface interface {
	CreateTodoList(ctx context.Context, input *input.CreateTodoListInput) (*output.CreateTodoListOutput, error)
	AddTodo(ctx context.Context, input *input.AddTodoInput) error
	GetTodoList(ctx context.Context, input *input.GetTodoListInput) (*output.GetTodoListOutput, error)
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

func (u *TodoUseCase) CreateTodoList(ctx context.Context, input *input.CreateTodoListInput) (*output.CreateTodoListOutput, error) {
	var result *output.CreateTodoListOutput
	err := u.tx.RWTx(ctx, func(ctx context.Context) error {
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

		todoList.MarkEventsAsCommitted()

		result = &output.CreateTodoListOutput{
			AggregateID: todoList.GetAggregateID().String(),
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *TodoUseCase) AddTodo(ctx context.Context, input *input.AddTodoInput) error {
	return u.addTodoWithRetry(ctx, input.AggregateID, input.UserID, input.Todo, 3)
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

			userIDVO, err := value.NewUserID(userID)
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

func (u *TodoUseCase) GetTodoList(ctx context.Context, input *input.GetTodoListInput) (*output.GetTodoListOutput, error) {
	var result *output.GetTodoListOutput

	err := u.tx.RWTx(ctx, func(ctx context.Context) error {
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

		items := make([]output.TodoItem, 0, len(todoList.GetItems()))
		for _, item := range todoList.GetItems() {
			items = append(items, output.TodoItem{
				Text: item.Text.String(),
			})
		}

		result = &output.GetTodoListOutput{
			AggregateID: todoList.GetAggregateID().String(),
			UserID:      todoList.GetUserID().String(),
			Items:       items,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

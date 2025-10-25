package usecase

import (
	"context"
	"fmt"

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

		// Save events to event store
		if err := u.eventStore.SaveEvents(ctx, todoList.GetAggregateID(), todoList.GetUncommittedEvents()); err != nil {
			return fmt.Errorf("failed to save events: %w", err)
		}
		
		todoList.MarkEventsAsCommitted()

		return nil
	})

	if err != nil {
		return  "", nil
	}

	return aggregateID, nil
}

func (u *TodoUseCase) AddTodo(ctx context.Context, aggregateID string, userID string, todo string) error {
	return u.tx.RWTx(ctx, func(ctx context.Context) error {
		todoText, err := value.NewTodoText(todo)
		if err != nil {
			return err
		}

		aggregateUUID, err := uuid.Parse(aggregateID)
		if err != nil {
			return fmt.Errorf("invalid aggregate ID: %w", err)
		}

		// TODO: 既存のTodoListを読み込み
		// todoList := eventStore.LoadAggregate(ctx, aggregateUUID)
		// if todoList == nil {
		//     return fmt.Errorf("todo list not found")
		// }
		todoList := aggregate.NewTodoListAggregate()
		// 暫定的に既存リストがあると仮定してaggregateIDを設定
		
		cmd := command.AddTodoCommand{
			AggregateID: aggregateUUID,
			UserID:      userID,
			TodoText:    todoText,
		}
		
		if err := todoList.ExecuteAddTodoCommand(cmd); err != nil {
			return fmt.Errorf("failed to handle add todo command: %w", err)
		}

		// Save events to event store
		if err := u.eventStore.SaveEvents(ctx, todoList.GetAggregateID(), todoList.GetUncommittedEvents()); err != nil {
			return fmt.Errorf("failed to save events: %w", err)
		}
		
		todoList.MarkEventsAsCommitted()

		return nil
	})
}

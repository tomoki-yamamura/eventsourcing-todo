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
	AddTodo(ctx context.Context, todo string) error
}

type TodoUseCase struct {
	tx repository.Transaction
}

func NewTodoUseCase(tx repository.Transaction) TodoUseCaseInterface {
	return &TodoUseCase{
		tx: tx,
	}
}

func (u *TodoUseCase) AddTodo(ctx context.Context, todo string) error {
	return u.tx.RWTx(ctx, func(ctx context.Context) error {
		aggregateID := uuid.New()
		todoText, err := value.NewTodoText(todo)
		if err != nil {
			return err
		}
		cmd := command.AddTodoCommand{
			AggregateID: aggregateID,
			TodoText:    todoText,
		}

		aggregate := aggregate.NewTodoAggregate()
		if err := aggregate.ExecuteAddTodoCommand(cmd); err != nil {
			return fmt.Errorf("failed to handle add todo command: %w", err)
		}

		aggregate.MarkEventsAsCommitted()

		return nil
	})
}

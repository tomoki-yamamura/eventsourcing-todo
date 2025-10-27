package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/aggregate"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/repository"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/projector"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/input"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/query/output"
)

type TodoQueryUseCaseInterface interface {
	GetTodoList(ctx context.Context, input *input.GetTodoListInput) (*output.GetTodoListOutput, error)
}

type TodoQueryUseCase struct {
	tx         repository.Transaction
	eventStore repository.EventStore
	projector  *projector.InMemTodoProjector
}

func NewTodoQueryUseCase(tx repository.Transaction, eventStore repository.EventStore, projector *projector.InMemTodoProjector) TodoQueryUseCaseInterface {
	return &TodoQueryUseCase{
		tx:         tx,
		eventStore: eventStore,
		projector:  projector,
	}
}

func (u *TodoQueryUseCase) GetTodoList(ctx context.Context, input *input.GetTodoListInput) (*output.GetTodoListOutput, error) {
	view := u.projector.GetList(input.AggregateID)
	if view != nil {
		items := make([]output.TodoItem, 0, len(view.Items))
		for _, item := range view.Items {
			items = append(items, output.TodoItem{
				Text: item.Text,
			})
		}

		return &output.GetTodoListOutput{
			AggregateID: view.AggregateID,
			UserID:      view.UserID,
			Items:       items,
		}, nil
	}

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

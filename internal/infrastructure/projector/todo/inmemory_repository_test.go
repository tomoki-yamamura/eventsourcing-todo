package todo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/errors"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/projector/todo"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/readmodelstore/dto"
)

func TestInMemoryTodoListViewRepository_Get(t *testing.T) {
	tests := map[string]struct {
		existingData map[string]*dto.TodoListViewDTO
		aggregateID  string
		want         *dto.TodoListViewDTO
		wantError    error
	}{
		"should return existing item": {
			existingData: map[string]*dto.TodoListViewDTO{
				"test-id": {
					AggregateID: "test-id",
					UserID:      "user123",
					Items: []dto.TodoItemViewDTO{
						{Text: "Task 1"},
					},
					Version:   1,
					UpdatedAt: time.Now(),
				},
			},
			aggregateID: "test-id",
			want: &dto.TodoListViewDTO{
				AggregateID: "test-id",
				UserID:      "user123",
				Items: []dto.TodoItemViewDTO{
					{Text: "Task 1"},
				},
				Version: 1,
			},
		},
		"should return error when aggregateID not found": {
			existingData: map[string]*dto.TodoListViewDTO{},
			aggregateID:  "non-existing",
			want:         nil,
			wantError:    errors.NotFound.New("todo list not found"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := todo.NewInMemoryTodoListViewRepository()

			// Pre-populate test data
			for id, view := range tt.existingData {
				err := repo.Save(context.Background(), id, view)
				require.NoError(t, err)
			}

			// Act
			got, err := repo.Get(context.Background(), tt.aggregateID)

			// Assert
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.AggregateID, got.AggregateID)
				require.Equal(t, tt.want.UserID, got.UserID)
				require.Equal(t, tt.want.Version, got.Version)
				require.Equal(t, len(tt.want.Items), len(got.Items))
				for i, item := range tt.want.Items {
					require.Equal(t, item.Text, got.Items[i].Text)
				}
			}
		})
	}
}

func TestInMemoryTodoListViewRepository_Save(t *testing.T) {
	tests := map[string]struct {
		aggregateID string
		view        *dto.TodoListViewDTO
		wantError   error
	}{
		"should save valid view": {
			aggregateID: "test-id",
			view: &dto.TodoListViewDTO{
				AggregateID: "test-id",
				UserID:      "user123",
				Items: []dto.TodoItemViewDTO{
					{Text: "Task 1"},
					{Text: "Task 2"},
				},
				Version:   1,
				UpdatedAt: time.Now(),
			},
		},
		"should handle nil view": {
			aggregateID: "test-id",
			view:        nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := todo.NewInMemoryTodoListViewRepository()

			// Act
			err := repo.Save(context.Background(), tt.aggregateID, tt.view)

			// Assert
			if tt.wantError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantError)
			} else {
				require.NoError(t, err)
				saved, err := repo.Get(context.Background(), tt.aggregateID)
				require.NoError(t, err)
				require.Equal(t, tt.view.AggregateID, saved.AggregateID)
				require.Equal(t, tt.view.UserID, saved.UserID)
				require.Equal(t, tt.view.Version, saved.Version)
			}
		})
	}
}

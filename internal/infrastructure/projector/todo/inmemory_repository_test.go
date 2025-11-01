package todo

import (
	"context"
	"maps"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/readmodelstore/dto"
)

func TestInMemoryTodoListViewRepository_Get(t *testing.T) {
	tests := map[string]struct {
		existingData map[string]*dto.TodoListViewDTO
		aggregateID  string
		want         *dto.TodoListViewDTO
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
		"should return nil when aggregateID not found": {
			existingData: map[string]*dto.TodoListViewDTO{},
			aggregateID:  "non-existing",
			want:         nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			repo := &InMemoryTodoListViewRepository{
				data: make(map[string]*dto.TodoListViewDTO),
			}
			maps.Copy(repo.data, tt.existingData)

			// Act
			got := repo.Get(context.Background(), tt.aggregateID)

			// Assert
			if tt.want == nil {
				require.Nil(t, got)
			} else {
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
			repo := &InMemoryTodoListViewRepository{
				data: make(map[string]*dto.TodoListViewDTO),
			}

			// Act
			err := repo.Save(context.Background(), tt.aggregateID, tt.view)

			// Assert
			require.NoError(t, err)

			if tt.view != nil {
				saved := repo.Get(context.Background(), tt.aggregateID)
				require.Equal(t, tt.view.AggregateID, saved.AggregateID)
				require.Equal(t, tt.view.UserID, saved.UserID)
				require.Equal(t, tt.view.Version, saved.Version)
			}
		})
	}
}

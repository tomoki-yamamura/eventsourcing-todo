package todo

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/bus"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/projector"
)

type TodoProjectorInterface interface {
	projector.Projector
	GetList(aggregateID string) *TodoListView
}

type TodoProjector struct {
	viewRepo TodoListViewRepository
	seen     map[string]struct{}
}

func NewTodoProjector(viewRepo TodoListViewRepository) TodoProjectorInterface {
	return &TodoProjector{
		viewRepo: viewRepo,
		seen:     make(map[string]struct{}),
	}
}

func (p *TodoProjector) Start(ctx context.Context, eventBus bus.EventBus) error {
	return eventBus.Subscribe(func(ctx context.Context, e event.Event) error {
		eventID := e.GetEventID().String()
		if _, ok := p.seen[eventID]; ok {
			return nil
		}
		p.seen[eventID] = struct{}{}

		aggregateID := e.GetAggregateID().String()
		currentView := p.viewRepo.Get(aggregateID)
		
		updatedView := p.applyToView(currentView, e)
		if updatedView != nil {
			return p.viewRepo.Save(aggregateID, updatedView)
		}
		
		return nil
	})
}

func (p *TodoProjector) GetList(aggregateID string) *TodoListView {
	return p.viewRepo.Get(aggregateID)
}

func (p *TodoProjector) applyToView(view *TodoListView, e event.Event) *TodoListView {
	switch evt := e.(type) {
	case event.TodoListCreatedEvent:
		return &TodoListView{
			AggregateID: evt.GetAggregateID().String(),
			UserID:      evt.UserID.String(),
			Items:       []TodoItemView{},
			Version:     evt.GetVersion(),
			UpdatedAt:   evt.GetTimestamp(),
		}
	case event.TodoAddedEvent:
		if view == nil {
			return nil // View not found, ignore
		}
		
		// Check version to ensure safe application
		if evt.GetVersion() <= view.Version {
			return view // Already applied or out of order
		}
		
		newItems := make([]TodoItemView, len(view.Items))
		copy(newItems, view.Items)
		newItems = append(newItems, TodoItemView{
			Text: evt.TodoText.String(),
		})
		
		return &TodoListView{
			AggregateID: view.AggregateID,
			UserID:      view.UserID,
			Items:       newItems,
			Version:     evt.GetVersion(),
			UpdatedAt:   evt.GetTimestamp(),
		}
	}
	
	return view
}
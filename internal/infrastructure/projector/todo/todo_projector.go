package todo

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports"
)

type TodoProjectorImpl struct {
	viewRepo *InMemoryTodoListViewRepository
	seen     map[string]struct{}
}

func NewTodoProjector(viewRepo *InMemoryTodoListViewRepository) ports.Projector {
	return &TodoProjectorImpl{
		viewRepo: viewRepo,
		seen:     make(map[string]struct{}),
	}
}

func (p *TodoProjectorImpl) Subscribe(handler func(context.Context, event.Event) error) error {
	return nil
}

func (p *TodoProjectorImpl) Start(ctx context.Context, eventBus ports.EventSubscriber) error {
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

func (p *TodoProjectorImpl) applyToView(view *TodoListViewDTO, e event.Event) *TodoListViewDTO {
	switch evt := e.(type) {
	case event.TodoListCreatedEvent:
		return &TodoListViewDTO{
			AggregateID: evt.GetAggregateID().String(),
			UserID:      evt.UserID.String(),
			Items:       []TodoItemViewDTO{},
			Version:     evt.GetVersion(),
			UpdatedAt:   evt.GetTimestamp(),
		}
	case event.TodoAddedEvent:
		if view == nil {
			return nil
		}

		if evt.GetVersion() <= view.Version {
			return view
		}

		newItems := make([]TodoItemViewDTO, len(view.Items))
		copy(newItems, view.Items)
		newItems = append(newItems, TodoItemViewDTO{
			Text: evt.TodoText.String(),
		})

		return &TodoListViewDTO{
			AggregateID: view.AggregateID,
			UserID:      view.UserID,
			Items:       newItems,
			Version:     evt.GetVersion(),
			UpdatedAt:   evt.GetTimestamp(),
		}
	}

	return view
}

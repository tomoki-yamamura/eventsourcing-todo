package todo

import (
	"context"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/gateway"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/readmodelstore"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/usecase/ports/readmodelstore/dto"
)

type TodoProjectorImpl struct {
	viewRepo readmodelstore.TodoListReadModelStore
	seen     map[string]struct{}
}

func NewTodoProjector(viewRepo readmodelstore.TodoListReadModelStore) gateway.Projector {
	return &TodoProjectorImpl{
		viewRepo: viewRepo,
		seen:     make(map[string]struct{}),
	}
}

func (p *TodoProjectorImpl) Handle(ctx context.Context, e event.Event) error {
	eventID := e.GetEventID().String()
	if _, ok := p.seen[eventID]; ok {
		return nil
	}
	p.seen[eventID] = struct{}{}

	aggID := e.GetAggregateID().String()
	current := p.viewRepo.Get(ctx, aggID)
	updated := p.applyToView(current, e)
	if updated != nil {
		return p.viewRepo.Save(ctx, aggID, updated)
	}
	return nil
}

func (p *TodoProjectorImpl) Start(ctx context.Context, bus gateway.EventSubscriber) error {
	bus.Subscribe(p.Handle)
	return nil
}

func (p *TodoProjectorImpl) applyToView(view *dto.TodoListViewDTO, e event.Event) *dto.TodoListViewDTO {
	switch evt := e.(type) {
	case event.TodoListCreatedEvent:
		return &dto.TodoListViewDTO{
			AggregateID: evt.GetAggregateID().String(),
			UserID:      evt.UserID.String(),
			Items:       []dto.TodoItemViewDTO{},
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

		newItems := make([]dto.TodoItemViewDTO, len(view.Items))
		copy(newItems, view.Items)
		newItems = append(newItems, dto.TodoItemViewDTO{
			Text: evt.TodoText.String(),
		})

		return &dto.TodoListViewDTO{
			AggregateID: view.AggregateID,
			UserID:      view.UserID,
			Items:       newItems,
			Version:     evt.GetVersion(),
			UpdatedAt:   evt.GetTimestamp(),
		}
	}

	return view
}

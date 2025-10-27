package projector

import (
	"context"
	"sync"
	"time"

	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/infrastructure/bus"
)

type TodoListView struct {
	AggregateID string
	UserID      string
	Items       []TodoItemView
	Version     int
	UpdatedAt   time.Time
}

type TodoItemView struct {
	Text string
}

type InMemTodoProjector struct {
	mu    sync.RWMutex
	views map[string]*TodoListView // by aggregateID
	seen  map[string]struct{}      // event_id deduplication
}

func NewInMemTodoProjector() *InMemTodoProjector {
	return &InMemTodoProjector{
		views: make(map[string]*TodoListView),
		seen:  make(map[string]struct{}),
	}
}

func (p *InMemTodoProjector) Start(ctx context.Context, eventBus bus.EventBus) error {
	return eventBus.Subscribe(func(ctx context.Context, e event.Event) error {
		p.mu.Lock()
		defer p.mu.Unlock()

		// Deduplication
		eventID := e.GetEventID().String()
		if _, ok := p.seen[eventID]; ok {
			return nil
		}
		p.seen[eventID] = struct{}{}

		aggregateID := e.GetAggregateID().String()
		view := p.views[aggregateID]
		
		// Apply event to view
		updatedView := p.applyToView(view, e)
		if updatedView != nil {
			p.views[aggregateID] = updatedView
		}
		
		return nil
	})
}

func (p *InMemTodoProjector) GetList(aggregateID string) *TodoListView {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	view := p.views[aggregateID]
	if view == nil {
		return nil
	}
	
	// Return a copy to avoid race conditions
	return p.cloneView(view)
}

func (p *InMemTodoProjector) applyToView(view *TodoListView, e event.Event) *TodoListView {
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

func (p *InMemTodoProjector) cloneView(view *TodoListView) *TodoListView {
	if view == nil {
		return nil
	}
	
	items := make([]TodoItemView, len(view.Items))
	copy(items, view.Items)
	
	return &TodoListView{
		AggregateID: view.AggregateID,
		UserID:      view.UserID,
		Items:       items,
		Version:     view.Version,
		UpdatedAt:   view.UpdatedAt,
	}
}
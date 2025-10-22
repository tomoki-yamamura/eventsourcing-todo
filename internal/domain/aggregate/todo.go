package aggregate

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/command"
	"github.com/tomoki-yamamura/eventsourcing-todo/internal/domain/event"
)

type TodoAggregate struct {
	aggregateID       uuid.UUID
	version           int
	uncommittedEvents []event.Event
}

func NewTodoAggregate() *TodoAggregate {
	return &TodoAggregate{
		uncommittedEvents: make([]event.Event, 0),
	}
}

func (a *TodoAggregate) GetAggregateID() uuid.UUID {
	return a.aggregateID
}

func (a *TodoAggregate) GetVersion() int {
	return a.version
}

func (a *TodoAggregate) GetUncommittedEvents() []event.Event {
	return a.uncommittedEvents
}

func (a *TodoAggregate) MarkEventsAsCommitted() {
	a.uncommittedEvents = make([]event.Event, 0)
}

func (a *TodoAggregate) Hydration(events []event.Event) error {
	for _, evt := range events {
		if err := a.applyEvent(evt, false); err != nil {
			return fmt.Errorf("failed to apply event: %w", err)
		}
	}
	return nil
}

func (a *TodoAggregate) ExecuteAddTodoCommand(cmd command.AddTodoCommand) error {
	if cmd.Todo == "" {
		return fmt.Errorf("todo cannot be empty")
	}

	evt := event.TodoAddedEvent{
		AggregateID: cmd.AggregateID,
		Todo:        cmd.Todo,
		EventID:     uuid.New(),
		Timestamp:   time.Now(),
		Version:     a.version + 1,
	}

	return a.applyEvent(evt, true)
}

func (a *TodoAggregate) applyEvent(evt event.Event, isNew bool) error {
	switch e := evt.(type) {
	case event.TodoAddedEvent:
		a.onTodoAdded(e)
	default:
		return fmt.Errorf("unknown event type: %T", evt)
	}

	if isNew {
		a.uncommittedEvents = append(a.uncommittedEvents, evt)
	}

	a.version = evt.GetVersion()
	return nil
}

func (a *TodoAggregate) onTodoAdded(evt event.TodoAddedEvent) {
	if a.aggregateID == uuid.Nil {
		a.aggregateID = evt.AggregateID
	}
}

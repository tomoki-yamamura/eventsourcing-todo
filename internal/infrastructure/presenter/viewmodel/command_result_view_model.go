package viewmodel

type CommandResultViewModel struct {
	AggregateID string           `json:"aggregateId"`
	Version     int              `json:"version"`
	Events      []EventViewModel `json:"events"`
	Status      string           `json:"status"`
	ExecutedAt  string           `json:"executedAt"`
}

type EventViewModel struct {
	Type       string `json:"type"`
	Version    int    `json:"version"`
	Data       any    `json:"data"`
	OccurredAt string `json:"occurredAt"`
}

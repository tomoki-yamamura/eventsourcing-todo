package viewmodel

type TodoListVM struct {
	AggregateID string     `json:"aggregate_id"`
	UserID      string     `json:"user_id"`
	Items       []TodoItem `json:"items"`
	UpdatedAt   string     `json:"updated_at"`
}

type TodoItem struct {
	Text string `json:"text"`
}

package output

type AddTodoOutput struct {
	AggregateID string     `json:"aggregate_id"`
	UserID      string     `json:"user_id"`
	Items       []TodoItem `json:"items"`
}

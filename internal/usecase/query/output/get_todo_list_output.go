package output

type GetTodoListOutput struct {
	AggregateID string
	UserID      string
	Items       []TodoItem
}

type TodoItem struct {
	Text string
}

package response

type CreateTodoListResponse struct {
	AggregateID string `json:"aggregate_id"`
	Message     string `json:"message"`
}

type AddTodoResponse struct {
	Message string `json:"message"`
}

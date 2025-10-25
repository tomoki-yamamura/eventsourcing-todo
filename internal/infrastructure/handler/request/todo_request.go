package request

type CreateTodoListRequest struct {
	UserID string `json:"user_id"`
}

type AddTodoRequest struct {
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

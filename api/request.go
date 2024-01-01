package api

type CreateTodoRequest struct {
	Contents string `json:"contents"`
}

type DeleteTodoRequest struct {
	ID int `json:"id"`
}

type EditTodoRequest struct {
	ID       int    `json:"id"`
	Contents string `json:"contents"`
}

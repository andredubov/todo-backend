package domain

type TodoList struct {
	Id          int    `json:"id,omitempty" db:"id"`
	Title       string `json:"title,omitempty" db:"title" validate:"nonzero"`
	Description string `json:"description,omitempty" db:"description"`
}

type UpdateTodoListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

package domain

type TodoItem struct {
	Id          int    `json:"id,omitempty" db:"id"`
	Title       string `json:"title,omitempty" db:"title" validate:"nonzero"`
	Description string `json:"description,omitempty" db:"description"`
	Done        bool   `json:"done,omitempty" db:"done"`
}

type UpdateTodoItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

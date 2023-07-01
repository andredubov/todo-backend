package domain

type User struct {
	Id       int    `json:"id,omitempty" db:"id"`
	Name     string `json:"name,omitempty" db:"name" validate:"min=3, max=40"`
	Email    string `json:"email,omitempty" db:"email" validate:"nonzero"`
	Password string `json:"password,omitempty" db:"password_hash" validate:"min=6"`
}

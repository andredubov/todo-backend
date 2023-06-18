package domain

type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name" validate:"min=3,max=40"`
	Email    string `json:"email" db:"email" validate:"nonzero"`
	Password string `json:"password" db:"password_hash" validate:"min=6"`
}

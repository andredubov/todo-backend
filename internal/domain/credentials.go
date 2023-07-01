package domain

type Credentials struct {
	Email    string `json:"email" validate:"nonzero"`
	Password string `json:"password" validate:"min=6"`
}

package repository

import (
	"context"
	"fmt"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

type postgresUsersRepository struct {
	db *sqlx.DB
}

func NewPostgresUsersRepository(db *sqlx.DB) *postgresUsersRepository {
	return &postgresUsersRepository{db: db}
}

func (r *postgresUsersRepository) Create(ctx context.Context, user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *postgresUsersRepository) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	todoListTable   = "todo_lists"
	usersListsTable = "users_lists"
)

type postgresTodoListRepository struct {
	db *sqlx.DB
}

func NewPostgresTodoListRepository(db *sqlx.DB) *postgresTodoListRepository {
	return &postgresTodoListRepository{db: db}
}

func (r *postgresTodoListRepository) Create(ctx context.Context, todolist domain.TodoList, userID int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var todoListId int
	createTodoListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListTable)
	row := r.db.QueryRow(createTodoListQuery, todolist.Title, todolist.Description)
	if err := row.Scan(&todoListId); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListsQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2) RETURNING id", usersListsTable)
	_, err = tx.Exec(createUsersListsQuery, userID, todoListId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return todoListId, tx.Commit()
}

func (r *postgresTodoListRepository) GetByUserId(ctx context.Context, userId int) ([]domain.TodoList, error) {

	var todolists []domain.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1", todoListTable, usersListsTable)
	err := r.db.Select(&todolists, query, userId)

	return todolists, err
}

func (r *postgresTodoListRepository) GetById(ctx context.Context, userId, listId int) (domain.TodoList, error) {

	var todolist domain.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2",
		todoListTable, usersListsTable)
	err := r.db.Get(&todolist, query, userId, listId)

	return todolist, err
}

func (r *postgresTodoListRepository) Delete(ctx context.Context, userId, listId int) error {

	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2", todoListTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *postgresTodoListRepository) Update(ctx context.Context, userId, listId int, input domain.UpdateTodoListInput) error {

	setValues, args, argId := make([]string, 0), make([]interface{}, 0), 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, listId, userId)

	_, err := r.db.Exec(query, args...)

	return err
}

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type postgresTodoItemRepository struct {
	db *sqlx.DB
}

func NewPostgresTodoItemRepository(db *sqlx.DB) *postgresTodoItemRepository {
	return &postgresTodoItemRepository{db: db}
}

func (r *postgresTodoItemRepository) Create(ctx context.Context, listId int, item domain.TodoItem) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *postgresTodoItemRepository) GetAll(ctx context.Context, userId, listId int) ([]domain.TodoItem, error) {
	var todoItems []domain.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
									INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`, todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Select(&todoItems, query, listId, userId); err != nil {
		return nil, err
	}

	return todoItems, nil
}

func (r *postgresTodoItemRepository) GetById(ctx context.Context, userId, itemId int) (domain.TodoItem, error) {
	var todoItem domain.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
									INNER JOIN %s li on li.item_id = ti.id
									INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`, todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Get(&todoItem, query, itemId, userId); err != nil {
		return todoItem, err
	}

	return todoItem, nil
}

func (r *postgresTodoItemRepository) Delete(ctx context.Context, userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, itemId)

	return err
}

func (r *postgresTodoItemRepository) Update(ctx context.Context, userId, itemId int, todoItem domain.TodoItem) error {

	setValues, args, argId := make([]string, 0), make([]interface{}, 0), 1

	setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
	args = append(args, todoItem.Title)
	argId++

	setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
	args = append(args, todoItem.Description)
	argId++

	setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
	args = append(args, todoItem.Done)
	argId++

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}

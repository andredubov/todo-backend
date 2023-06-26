package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/dvln/testify/assert"
	"github.com/jmoiron/sqlx"
)

func TestItem_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	todoItemRepository := NewPostgresTodoItemRepository(dbx)

	type (
		args struct {
			listId int
			item   domain.TodoItem
		}

		mockBehavior func(args args, id int)

		test struct {
			name         string
			input        args
			mockBehavior mockBehavior
			wantId       int
			wantErr      bool
		}
	)

	tests := []test{
		{
			name: "Ok",
			input: args{
				listId: 1,
				item: domain.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				itemsTableQuery, listsItemsTableQuery := fmt.Sprintf("INSERT INTO %s", todoItemsTable), fmt.Sprintf("INSERT INTO %s", listsItemsTable)
				mock.ExpectQuery(itemsTableQuery).WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)
				mock.ExpectExec(listsItemsTableQuery).WithArgs(args.listId, id).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantId:  1,
			wantErr: false,
		},
		{
			name: "Empty Fields",
			input: args{
				listId: 1,
				item: domain.TodoItem{
					Title:       "",
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("insert error"))
				itemsTableQuery := fmt.Sprintf("INSERT INTO %s", todoItemsTable)
				mock.ExpectQuery(itemsTableQuery).WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failed second insert",
			input: args{
				listId: 1,
				item: domain.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				itemsTableQuery, listsItemsTableQuery := fmt.Sprintf("INSERT INTO %s", todoItemsTable), fmt.Sprintf("INSERT INTO %s", listsItemsTable)
				mock.ExpectQuery(itemsTableQuery).WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)
				mock.ExpectExec(listsItemsTableQuery).WithArgs(args.listId, id).WillReturnError(errors.New("insert error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mockBehavior(test.input, test.wantId)

			gotId, err := todoItemRepository.Create(context.TODO(), test.input.listId, test.input.item)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantId, gotId)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItem_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	todoItemRepository := NewPostgresTodoItemRepository(dbx)

	type (
		args struct {
			listId int
			userId int
		}
		test struct {
			name         string
			mockBehavior func()
			input        args
			want         []domain.TodoItem
			wantErr      bool
		}
	)

	tests := []test{
		{
			name: "Ok",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title1", "description1", true).
					AddRow(2, "title2", "description2", false).
					AddRow(3, "title3", "description3", false)

				query := fmt.Sprintf("SELECT (.+) FROM %s ti INNER JOIN %s li on (.+) INNER JOIN %s ul on (.+) WHERE (.+)", todoItemsTable, listsItemsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
			want: []domain.TodoItem{
				{Id: 1, Title: "title1", Description: "description1", Done: true},
				{Id: 2, Title: "title2", Description: "description2", Done: false},
				{Id: 3, Title: "title3", Description: "description3", Done: false},
			},
		},
		{
			name: "No Records",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})
				query := fmt.Sprintf("SELECT (.+) FROM %s ti INNER JOIN %s li on (.+) INNER JOIN %s ul on (.+) WHERE (.+)", todoItemsTable, listsItemsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mockBehavior()

			got, err := todoItemRepository.GetAll(context.TODO(), test.input.userId, test.input.listId)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

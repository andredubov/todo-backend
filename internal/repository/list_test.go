package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/dvln/testify/assert"
	"github.com/jmoiron/sqlx"
)

func TestList_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	todoListRepository := NewPostgresTodoListRepository(dbx)

	type (
		args struct {
			userId   int
			todoList domain.TodoList
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
				userId: 1,
				todoList: domain.TodoList{
					Title:       "title",
					Description: "description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				itemsTableQuery, usersListsTableQuery := fmt.Sprintf("INSERT INTO %s", todoListTable), fmt.Sprintf("INSERT INTO %s", usersListsTable)
				mock.ExpectQuery(itemsTableQuery).WithArgs(args.todoList.Title, args.todoList.Description).WillReturnRows(rows)
				mock.ExpectExec(usersListsTableQuery).WithArgs(args.userId, id).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantId:  1,
			wantErr: false,
		},
		{
			name: "Empty fields",
			input: args{
				userId: 1,
				todoList: domain.TodoList{
					Description: "description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				itemsTableQuery := fmt.Sprintf("INSERT INTO %s", todoListTable)
				mock.ExpectQuery(itemsTableQuery).WithArgs(args.todoList.Title, args.todoList.Description).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			wantId:  1,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mockBehavior(test.input, test.wantId)

			gotId, err := todoListRepository.Create(context.TODO(), test.input.todoList, test.input.userId)
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

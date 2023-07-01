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

func TestUser_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	usersRepository := NewPostgresUsersRepository(dbx)

	type (
		test struct {
			name         string
			mockBehavior func(domain.User)
			input        domain.User
			want         int
			wantErr      bool
		}
	)

	tests := []test{
		{
			name: "Ok",
			mockBehavior: func(user domain.User) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				query := fmt.Sprintf("INSERT INTO %s", usersTable)
				mock.ExpectQuery(query).WithArgs(user.Name, user.Email, user.Password).WillReturnRows(rows)
			},
			input: domain.User{
				Name:     "test name",
				Email:    "test email",
				Password: "test password",
			},
			want: 1,
		},
		{
			name: "Empty fields",
			mockBehavior: func(user domain.User) {
				rows := sqlmock.NewRows([]string{"id"})
				query := fmt.Sprintf("INSERT INTO %s", usersTable)
				mock.ExpectQuery(query).WithArgs(user.Name, user.Email, user.Password).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mockBehavior(test.input)

			got, err := usersRepository.Create(context.TODO(), test.input)
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

func TestUser_GetByCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")
	usersRepository := NewPostgresUsersRepository(dbx)

	type (
		test struct {
			name         string
			mockBehavior func(domain.Credentials, domain.User)
			input        domain.Credentials
			want         domain.User
			wantErr      bool
		}
	)

	tests := []test{
		{
			name: "Ok",
			mockBehavior: func(cred domain.Credentials, user domain.User) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash"}).AddRow(user.Id, user.Name, user.Email, user.Password)
				query := fmt.Sprintf("SELECT (.+) FROM %s", usersTable)
				mock.ExpectQuery(query).WithArgs(cred.Email, cred.Password).WillReturnRows(rows)
			},
			input: domain.Credentials{
				Email:    "test email",
				Password: "test password",
			},
			want: domain.User{
				Id:       1,
				Name:     "test name",
				Email:    "test email",
				Password: "test password",
			},
		},
		{
			name: "Not found",
			mockBehavior: func(cred domain.Credentials, user domain.User) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash"})
				query := fmt.Sprintf("SELECT (.+) FROM %s", usersTable)
				mock.ExpectQuery(query).WithArgs(cred.Email, cred.Password).WillReturnRows(rows)
			},
			input: domain.Credentials{
				Email:    "email not found",
				Password: "password not found",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mockBehavior(test.input, test.want)

			got, err := usersRepository.GetByCredentials(context.TODO(), test.input.Email, test.input.Password)
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

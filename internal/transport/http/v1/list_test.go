package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/andredubov/todo-backend/internal/config"
	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/service"
	mock_service "github.com/andredubov/todo-backend/internal/service/mocks"
	"github.com/andredubov/todo-backend/pkg/auth"
	"github.com/dvln/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestHandler_createList(t *testing.T) {

	type (
		enviroment struct {
			appEnv               string
			httpHost             string
			httpPort             string
			postgresHost         string
			postgresPort         string
			postgresDatabaseName string
			postgresUsername     string
			postgresPassword     string
			postgressSSLMode     string
			passwordSalt         string
			jwtSigningKey        string
		}

		args struct {
			userId   int
			todoList domain.TodoList
		}

		mockBehavior func(s *mock_service.MockTodoList, args args)

		test struct {
			enviroment           enviroment
			name                 string
			jwtTTL               time.Duration
			delay                time.Duration
			inputRequestBody     string
			input                args
			mockBehavior         mockBehavior
			expectedStatusCode   int
			expectedResponseBody string
		}
	)

	setEnv := func(env enviroment) {
		os.Setenv(config.ApplicationEnvironment, env.appEnv)
		os.Setenv(config.HttpHost, env.httpHost)
		os.Setenv(config.HttpPort, env.httpPort)
		os.Setenv(config.PostgresHost, env.postgresHost)
		os.Setenv(config.PostgresPort, env.postgresPort)
		os.Setenv(config.PostgresDatabaseName, env.postgresDatabaseName)
		os.Setenv(config.PostgresUsername, env.postgresUsername)
		os.Setenv(config.PostgresPassword, env.postgresPassword)
		os.Setenv(config.PostgresSSLMode, env.postgressSSLMode)
		os.Setenv(config.PasswordSalt, env.passwordSalt)
		os.Setenv(config.JwtSigningKey, env.jwtSigningKey)
	}

	tests := []test{
		{
			enviroment: enviroment{
				appEnv:               "local",
				httpHost:             "localhost",
				httpPort:             "8080",
				postgresHost:         "localhost",
				postgresPort:         "5432",
				postgresDatabaseName: "postgres",
				postgresUsername:     "postgres",
				postgresPassword:     "qwerty",
				postgressSSLMode:     "disable",
				passwordSalt:         "salt",
				jwtSigningKey:        "key",
			},
			name:             "OK",
			jwtTTL:           time.Duration(5 * time.Minute),
			delay:            time.Duration(0 * time.Millisecond),
			inputRequestBody: `{"title": "test title", "description": "test description"}`,
			input: args{
				userId:   1,
				todoList: domain.TodoList{Title: "test title", Description: "test description"},
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {
				gomock.InOrder(
					s.EXPECT().Validate(args.todoList).Return(nil),
					s.EXPECT().Create(gomock.Any(), args.todoList, args.userId).Return(1, nil),
				)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1}\n",
		},
		{
			enviroment: enviroment{
				appEnv:               "local",
				httpHost:             "localhost",
				httpPort:             "8080",
				postgresHost:         "localhost",
				postgresPort:         "5432",
				postgresDatabaseName: "postgres",
				postgresUsername:     "postgres",
				postgresPassword:     "qwerty",
				postgressSSLMode:     "disable",
				passwordSalt:         "salt",
				jwtSigningKey:        "key",
			},
			name:             "No Title",
			jwtTTL:           time.Duration(5 * time.Minute),
			delay:            time.Duration(0 * time.Millisecond),
			inputRequestBody: `{"description": "test description"}`,
			input: args{
				userId:   1,
				todoList: domain.TodoList{Description: "test description"},
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {
				gomock.InOrder(
					s.EXPECT().Validate(args.todoList).Return(errors.New("the given data was not valid: Title: ")),
				)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\": \"the given data was not valid: Title: \"}",
		},
		{
			enviroment: enviroment{
				appEnv:               "local",
				httpHost:             "localhost",
				httpPort:             "8080",
				postgresHost:         "localhost",
				postgresPort:         "5432",
				postgresDatabaseName: "postgres",
				postgresUsername:     "postgres",
				postgresPassword:     "qwerty",
				postgressSSLMode:     "disable",
				passwordSalt:         "salt",
				jwtSigningKey:        "key",
			},
			name:             "No Description",
			jwtTTL:           time.Duration(5 * time.Minute),
			delay:            time.Duration(0 * time.Millisecond),
			inputRequestBody: `{"title": "test title"}`,
			input: args{
				userId:   1,
				todoList: domain.TodoList{Title: "test title"},
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {
				gomock.InOrder(
					s.EXPECT().Validate(args.todoList).Return(nil),
					s.EXPECT().Create(gomock.Any(), args.todoList, args.userId).Return(1, nil),
				)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1}\n",
		},
		{
			enviroment: enviroment{
				appEnv:               "local",
				httpHost:             "localhost",
				httpPort:             "8080",
				postgresHost:         "localhost",
				postgresPort:         "5432",
				postgresDatabaseName: "postgres",
				postgresUsername:     "postgres",
				postgresPassword:     "qwerty",
				postgressSSLMode:     "disable",
				passwordSalt:         "salt",
				jwtSigningKey:        "key",
			},
			name:             "Token Expired",
			jwtTTL:           time.Duration(1 * time.Millisecond),
			delay:            time.Duration(1 * time.Second),
			inputRequestBody: `{"title": "test title", "description": "test description"}`,
			input: args{
				userId:   1,
				todoList: domain.TodoList{Title: "test title", Description: "test description"},
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {

			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "{\"message\": \"Token is expired\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockTodoListService := mock_service.NewMockTodoList(controller)
			test.mockBehavior(mockTodoListService, test.input)

			setEnv(test.enviroment)

			cfg, err := config.Init(configPath)
			if err != nil {
				t.Errorf("config initializing failed: %s", err.Error())
				return
			}

			tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
			if err != nil {
				t.Error(err)
				return
			}

			token, err := tokenManager.NewJWT(strconv.Itoa(test.input.userId), test.jwtTTL)
			if err != nil {
				t.Error(err)
				return
			}

			<-time.After(test.delay)

			services := service.Service{TodoList: mockTodoListService}
			h := NewHandler(&services, tokenManager, cfg.Auth.JWT)

			router := mux.NewRouter()
			postRouter := router.Methods(http.MethodPost).Subrouter()
			postRouter.HandleFunc("/api/lists", h.createList)
			postRouter.Use(h.userIdentity)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/lists", bytes.NewBufferString(test.inputRequestBody))
			r.Header.Set(authorizationHeader, bearer+" "+token) // set jwt token
			router.ServeHTTP(w, r)                              // perforn request

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

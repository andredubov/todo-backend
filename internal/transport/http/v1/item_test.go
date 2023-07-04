package handler

import (
	"bytes"
	"errors"
	"fmt"
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

func TestHandler_createItem(t *testing.T) {

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
			userId     int
			todoListId int
			todoItem   domain.TodoItem
		}

		mockBehavior func(s *mock_service.MockTodoItem, args args)

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
			inputRequestBody: `{"title": "test title", "description": "test description", "done": false}`,
			input: args{
				userId:     1,
				todoListId: 1,
				todoItem:   domain.TodoItem{Title: "test title", Description: "test description", Done: false},
			},
			mockBehavior: func(s *mock_service.MockTodoItem, args args) {
				gomock.InOrder(
					s.EXPECT().Validate(args.todoItem).Return(nil),
					s.EXPECT().Create(gomock.Any(), args.todoListId, args.todoItem).Return(1, nil),
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
			inputRequestBody: `{"description": "test description", "done": false}`,
			input: args{
				userId:     1,
				todoListId: 2,
				todoItem:   domain.TodoItem{Description: "test description", Done: false},
			},
			mockBehavior: func(s *mock_service.MockTodoItem, args args) {
				gomock.InOrder(
					s.EXPECT().Validate(args.todoItem).Return(errors.New("the given data was not valid: Title: ")),
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
			inputRequestBody: `{"title": "test title", "done": false}`,
			input: args{
				userId:     1,
				todoListId: 3,
				todoItem:   domain.TodoItem{Title: "test title", Done: false},
			},
			mockBehavior: func(s *mock_service.MockTodoItem, args args) {
				gomock.InOrder(
					s.EXPECT().Validate(args.todoItem).Return(nil),
					s.EXPECT().Create(gomock.Any(), args.todoListId, args.todoItem).Return(1, nil),
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
			name:             "No Done",
			jwtTTL:           time.Duration(5 * time.Minute),
			delay:            time.Duration(0 * time.Millisecond),
			inputRequestBody: `{"title": "test title", "description": "test description"}`,
			input: args{
				userId:     1,
				todoListId: 4,
				todoItem:   domain.TodoItem{Title: "test title", Description: "test description"},
			},
			mockBehavior: func(s *mock_service.MockTodoItem, args args) {
				gomock.InOrder(
					s.EXPECT().Validate(args.todoItem).Return(nil),
					s.EXPECT().Create(gomock.Any(), args.todoListId, args.todoItem).Return(1, nil),
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
				userId:     1,
				todoListId: 1,
				todoItem:   domain.TodoItem{Title: "test title", Description: "test description"},
			},
			mockBehavior: func(s *mock_service.MockTodoItem, args args) {

			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "{\"message\": \"Token is expired\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockTodoItemService := mock_service.NewMockTodoItem(controller)
			test.mockBehavior(mockTodoItemService, test.input)

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

			services := service.Service{TodoItem: mockTodoItemService}
			h := NewHandler(&services, tokenManager, cfg.Auth.JWT)

			router := mux.NewRouter()
			postRouter := router.Methods(http.MethodPost).Subrouter()
			postRouter.HandleFunc("/api/lists/{id:[0-9]+}/items", h.createItem)
			postRouter.Use(h.userIdentity)

			w, endpoint := httptest.NewRecorder(), fmt.Sprintf("/api/lists/%d/items", test.input.todoListId)
			r := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(test.inputRequestBody))
			r.Header.Set(authorizationHeader, bearer+" "+token) // set jwt token
			router.ServeHTTP(w, r)                              // perforn request

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getItems(t *testing.T) {

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
			userId     int
			todoListId int
		}

		mockBehavior func(s *mock_service.MockTodoItem, args args)

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
			name:   "OK",
			jwtTTL: time.Duration(5 * time.Minute),
			delay:  time.Duration(0 * time.Millisecond),
			input: args{
				userId:     1,
				todoListId: 1,
			},
			mockBehavior: func(s *mock_service.MockTodoItem, args args) {
				todoItems := []domain.TodoItem{
					{Id: 1, Title: "title1", Description: "description1", Done: true},
				}
				s.EXPECT().GetAll(gomock.Any(), args.userId, args.todoListId).Return(todoItems, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"data\":[{\"id\":1,\"title\":\"title1\",\"description\":\"description1\",\"done\":true}]}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockTodoItemService := mock_service.NewMockTodoItem(controller)
			test.mockBehavior(mockTodoItemService, test.input)

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

			services := service.Service{TodoItem: mockTodoItemService}
			h := NewHandler(&services, tokenManager, cfg.Auth.JWT)

			router := mux.NewRouter()
			getRouter := router.Methods(http.MethodGet).Subrouter()
			getRouter.HandleFunc("/api/lists/{id:[0-9]+}/items", h.getItems)
			getRouter.Use(h.userIdentity)

			w, endpoint := httptest.NewRecorder(), fmt.Sprintf("/api/lists/%d/items", test.input.todoListId)
			r := httptest.NewRequest(http.MethodGet, endpoint, bytes.NewBufferString(test.inputRequestBody))
			r.Header.Set(authorizationHeader, bearer+" "+token) // set jwt token
			router.ServeHTTP(w, r)                              // perforn request

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

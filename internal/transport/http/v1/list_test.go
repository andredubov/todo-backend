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

func TestHandler_getLists(t *testing.T) {

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
			userId int
		}

		mockBehavior func(s *mock_service.MockTodoList, args args)

		test struct {
			enviroment           enviroment
			jwtTTL               time.Duration
			delay                time.Duration
			name                 string
			mockBehavior         mockBehavior
			input                args
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
				userId: 1,
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {
				todoLists := []domain.TodoList{
					{Id: 1, Title: "title1", Description: "description1"},
				}
				s.EXPECT().GetByUserId(gomock.Any(), args.userId).Return(todoLists, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"data\":[{\"id\":1,\"title\":\"title1\",\"description\":\"description1\"}]}\n",
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
			name:   "Empty",
			jwtTTL: time.Duration(5 * time.Minute),
			delay:  time.Duration(0 * time.Millisecond),
			input: args{
				userId: 1,
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {
				todoLists := []domain.TodoList{}
				s.EXPECT().GetByUserId(gomock.Any(), args.userId).Return(todoLists, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"data\":[]}\n",
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
			name:   "Expired Token",
			jwtTTL: time.Duration(5 * time.Millisecond),
			delay:  time.Duration(1 * time.Second),
			input: args{
				userId: 1,
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
			getRouter := router.Methods(http.MethodGet).Subrouter()
			getRouter.HandleFunc("/api/lists", h.getLists)
			getRouter.Use(h.userIdentity)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/lists", bytes.NewBufferString(""))
			r.Header.Set(authorizationHeader, bearer+" "+token) // set jwt token
			router.ServeHTTP(w, r)                              // perforn request

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getListByID(t *testing.T) {

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

		mockBehavior func(s *mock_service.MockTodoList, args args)

		test struct {
			enviroment           enviroment
			jwtTTL               time.Duration
			delay                time.Duration
			name                 string
			mockBehavior         mockBehavior
			input                args
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
				todoListId: 2,
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {
				todoList := domain.TodoList{Id: 1, Title: "title1", Description: "description1"}
				s.EXPECT().GetById(gomock.Any(), args.userId, args.todoListId).Return(todoList, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1,\"title\":\"title1\",\"description\":\"description1\"}\n",
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
			name:   "Empty",
			jwtTTL: time.Duration(5 * time.Minute),
			delay:  time.Duration(0 * time.Millisecond),
			input: args{
				userId:     1,
				todoListId: 2,
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {
				s.EXPECT().GetById(gomock.Any(), args.userId, args.todoListId).Return(domain.TodoList{}, errors.New(""))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "{\"message\": \"unable to get a todolist by id: \"}",
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
			name:   "Token Expired",
			jwtTTL: time.Duration(5 * time.Millisecond),
			delay:  time.Duration(1 * time.Second),
			input: args{
				userId:     1,
				todoListId: 2,
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
			getRouter := router.Methods(http.MethodGet).Subrouter()
			getRouter.HandleFunc("/api/lists/{id:[0-9]+}", h.getListByID)
			getRouter.Use(h.userIdentity)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/lists/%d", test.input.todoListId), bytes.NewBufferString(""))
			r.Header.Set(authorizationHeader, bearer+" "+token) // set jwt token
			router.ServeHTTP(w, r)                              // perforn request

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_updateListByID(t *testing.T) {

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
			userId              int
			todoListId          int
			updateTodoListInput domain.UpdateTodoListInput
		}

		mockBehavior func(s *mock_service.MockTodoList, args args)

		test struct {
			enviroment           enviroment
			jwtTTL               time.Duration
			delay                time.Duration
			name                 string
			mockBehavior         mockBehavior
			input                args
			inputRequestBody     string
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
			inputRequestBody: `{"title": "new title", "description": "new description"}`,
			input: args{
				userId:     1,
				todoListId: 2,
				updateTodoListInput: domain.UpdateTodoListInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
				},
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {
				s.EXPECT().Update(gomock.Any(), args.userId, args.todoListId, args.updateTodoListInput).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"status\":\"success\"}\n",
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
			inputRequestBody: `{"description": "new description"}`,
			input: args{
				userId:     1,
				todoListId: 2,
				updateTodoListInput: domain.UpdateTodoListInput{
					Description: stringPointer("new description"),
				},
			},
			mockBehavior: func(s *mock_service.MockTodoList, args args) {
				s.EXPECT().Update(gomock.Any(), args.userId, args.todoListId, args.updateTodoListInput).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"status\":\"success\"}\n",
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
			putRouter := router.Methods(http.MethodPut).Subrouter()
			putRouter.HandleFunc("/api/lists/{id:[0-9]+}", h.updateListByID)
			putRouter.Use(h.userIdentity)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/lists/%d", test.input.todoListId), bytes.NewBufferString(test.inputRequestBody))
			r.Header.Set(authorizationHeader, bearer+" "+token) // set jwt token
			router.ServeHTTP(w, r)                              // perforn request

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

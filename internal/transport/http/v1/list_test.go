package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/andredubov/todo-backend/internal/config"
	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/service"
	mock_service "github.com/andredubov/todo-backend/internal/service/mocks"
	mock_auth "github.com/andredubov/todo-backend/pkg/auth/mocks"
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
			jwtToken string
		}

		mockBehavior func(s *mock_service.MockTodoList, m *mock_auth.MockTokenManager, args args)

		test struct {
			enviroment           enviroment
			name                 string
			httpHeaderName       string
			httpHeaderValue      string
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
			httpHeaderName:   authorizationHeader,
			httpHeaderValue:  bearer + " token",
			inputRequestBody: `{"title": "test title", "description": "test description"}`,
			input: args{
				userId:   1,
				todoList: domain.TodoList{Title: "test title", Description: "test description"},
				jwtToken: "token",
			},
			mockBehavior: func(s *mock_service.MockTodoList, m *mock_auth.MockTokenManager, args args) {
				gomock.InOrder(
					m.EXPECT().Parse(args.jwtToken).Return(strconv.Itoa(args.userId), nil),
					s.EXPECT().Validate(args.todoList).Return(nil),
					s.EXPECT().Create(gomock.Any(), args.todoList, args.userId).Return(1, nil),
				)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockTodoListService := mock_service.NewMockTodoList(controller)
			mockTokenManager := mock_auth.NewMockTokenManager(controller)
			test.mockBehavior(mockTodoListService, mockTokenManager, test.input)

			setEnv(test.enviroment)

			cfg, err := config.Init(configPath)
			if err != nil {
				t.Errorf("config initializing failed: %s", err.Error())
				return
			}

			services := service.Service{TodoList: mockTodoListService}
			h := NewHandler(&services, mockTokenManager, cfg.Auth.JWT)

			// test server
			router := mux.NewRouter()
			postRouter := router.Methods(http.MethodPost).Subrouter()
			postRouter.HandleFunc("/api/lists", h.createList)
			postRouter.Use(h.userIdentity)

			// test
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/lists", bytes.NewBufferString(test.inputRequestBody))
			r.Header.Set(test.httpHeaderName, test.httpHeaderValue)
			router.ServeHTTP(w, r) // perforn request

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

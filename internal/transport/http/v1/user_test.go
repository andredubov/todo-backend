package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andredubov/todo-backend/internal/config"
	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/service"
	mock_service "github.com/andredubov/todo-backend/internal/service/mocks"
	"github.com/andredubov/todo-backend/pkg/auth"
	"github.com/dvln/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

const (
	configPath = "../../../../configs"
)

func TestHandler_signUp(t *testing.T) {

	type mockBehavior func(s *mock_service.MockUsers, user domain.User)

	type enviroment struct {
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

	type testCase struct {
		enviroment           enviroment
		name                 string
		inputRequestBody     string
		inputUser            domain.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}

	testCases := []testCase{
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
			inputRequestBody: `{"name": "Alex", "email": "alex@gmail.com", "password": "qwerty"}`,
			inputUser:        domain.User{Name: "Alex", Email: "alex@gmail.com", Password: "qwerty"},
			mockBehavior: func(s *mock_service.MockUsers, user domain.User) {
				first := s.EXPECT().Validate(user).Return(nil)
				s.EXPECT().Create(gomock.Any(), user).Return(1, nil).After(first)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1,\"name\":\"Alex\",\"email\":\"alex@gmail.com\",\"password\":\"qwerty\"}\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockUsersService := mock_service.NewMockUsers(controller)
			testCase.mockBehavior(mockUsersService, testCase.inputUser)

			setEnv(testCase.enviroment)

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

			services := service.Service{Users: mockUsersService}
			h := NewHandler(&services, tokenManager, cfg.Auth.JWT)

			// test server
			router := mux.NewRouter()
			authRouter := router.Methods(http.MethodPost).Subrouter()
			authRouter.HandleFunc("/auth/sign-up", h.signUp)

			// test
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(testCase.inputRequestBody))
			router.ServeHTTP(w, r) // perforn request

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

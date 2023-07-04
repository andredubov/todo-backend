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
	mock_auth "github.com/andredubov/todo-backend/pkg/auth/mocks"
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
			name:             "Username is less than min available length",
			inputRequestBody: "{\"name\": \"ad\", \"email\": \"alex@gmail.com\", \"password\": \"qwerty\"}",
			inputUser:        domain.User{Name: "ad", Email: "alex@gmail.com", Password: "qwerty"},
			mockBehavior: func(s *mock_service.MockUsers, user domain.User) {
				s.EXPECT().Validate(user).Return(errors.New("Name: less than min"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\": \"Name: less than min\"}",
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
			name:             "Username is greater than max available length",
			inputRequestBody: "{\"name\": \"12345678912345678912345678912345678912345\", \"email\": \"alex@gmail.com\", \"password\": \"qwerty\"}",
			inputUser:        domain.User{Name: "12345678912345678912345678912345678912345", Email: "alex@gmail.com", Password: "qwerty"},
			mockBehavior: func(s *mock_service.MockUsers, user domain.User) {
				s.EXPECT().Validate(user).Return(errors.New("Name: greater than max"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\": \"Name: greater than max\"}",
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
			name:             "Password is less than min available length",
			inputRequestBody: "{\"name\": \"Alex\", \"email\": \"alex@gmail.com\", \"password\": \"qwert\"}",
			inputUser:        domain.User{Name: "Alex", Email: "alex@gmail.com", Password: "qwert"},
			mockBehavior: func(s *mock_service.MockUsers, user domain.User) {
				s.EXPECT().Validate(user).Return(errors.New("Password: less than min"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\": \"Password: less than min\"}",
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
			name:             "Incorrect email",
			inputRequestBody: "{\"name\": \"Alex\", \"email\": \"alexgmailcom\", \"password\": \"qwerty\"}",
			inputUser:        domain.User{Name: "Alex", Email: "alexgmailcom", Password: "qwerty"},
			mockBehavior: func(s *mock_service.MockUsers, user domain.User) {
				s.EXPECT().Validate(user).Return(errors.New("mail: missing '@' or angle-addr"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"message\": \"mail: missing '@' or angle-addr\"}",
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

func TestHandler_signIn(t *testing.T) {

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
			userId      int
			credentials domain.Credentials
			jwtCfg      config.JWTConfig
		}

		mockBehavior func(s *mock_service.MockUsers, m *mock_auth.MockTokenManager, args args, output SignInResponse)

		test struct {
			enviroment           enviroment
			name                 string
			mockBehavior         mockBehavior
			input                args
			output               SignInResponse
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
			inputRequestBody: `{"email": "user@gmail.com", "password": "qwerty"}`,
			input: args{
				userId:      1,
				credentials: domain.Credentials{Email: "user@gmail.com", Password: "qwerty"},
				jwtCfg: config.JWTConfig{
					AccessTokenTTL:  5 * time.Minute,
					RefreshTokenTTL: 5 * time.Minute,
					SigningKey:      "sign",
				},
			},
			output: SignInResponse{
				AccessToken:   "accessToken",
				ResfreshToken: "refreshToken",
			},
			mockBehavior: func(s *mock_service.MockUsers, m *mock_auth.MockTokenManager, input args, output SignInResponse) {
				userId := strconv.Itoa(input.userId)
				gomock.InOrder(
					s.EXPECT().GetByCredentials(gomock.Any(), input.credentials).Return(domain.User{Id: input.userId}, nil),
					m.EXPECT().NewJWT(userId, input.jwtCfg.AccessTokenTTL).Return(output.AccessToken, nil),
					m.EXPECT().NewJWT(userId, input.jwtCfg.RefreshTokenTTL).Return(output.ResfreshToken, nil),
				)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"accessToken\":\"accessToken\",\"refreshToken\":\"refreshToken\"}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockUsersService := mock_service.NewMockUsers(controller)
			mockTokenManger := mock_auth.NewMockTokenManager(controller)
			test.mockBehavior(mockUsersService, mockTokenManger, test.input, test.output)

			setEnv(test.enviroment)

			services := service.Service{Users: mockUsersService}
			h := NewHandler(&services, mockTokenManger, test.input.jwtCfg)

			router := mux.NewRouter()
			postRouter := router.Methods(http.MethodPut).Subrouter()
			postRouter.HandleFunc("/auth/sign-in", h.signIn)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/auth/sign-in", bytes.NewBufferString(test.inputRequestBody))
			router.ServeHTTP(w, r) // perforn request

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

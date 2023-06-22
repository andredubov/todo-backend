package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/andredubov/todo-backend/internal/config"
	"github.com/andredubov/todo-backend/internal/service"
	mock_auth "github.com/andredubov/todo-backend/pkg/auth/mocks"
	"github.com/dvln/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestHandler_userIdentity(t *testing.T) {

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

	type mockBehavior func(s *mock_auth.MockTokenManager, token string)

	type testCase struct {
		enviroment           enviroment
		name                 string
		headerName           string
		headerValue          string
		token                string
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
			name:        "OK",
			headerName:  authorizationHeader,
			headerValue: bearer + " token",
			token:       "token",
			mockBehavior: func(m *mock_auth.MockTokenManager, token string) {
				m.EXPECT().Parse(token).Return("1", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "1",
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
			name:                 "No header",
			headerName:           "",
			headerValue:          "",
			token:                "",
			mockBehavior:         func(m *mock_auth.MockTokenManager, token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: "{\"message\": \"empty auth header\"}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			mockTokenManager := mock_auth.NewMockTokenManager(controller)
			testCase.mockBehavior(mockTokenManager, testCase.token)

			setEnv(testCase.enviroment)

			cfg, err := config.Init(configPath)
			if err != nil {
				t.Errorf("config initializing failed: %s", err.Error())
				return
			}

			h := NewHandler(&service.Service{}, mockTokenManager, cfg.Auth.JWT)

			// test server
			router := mux.NewRouter()
			getRouter := router.Methods(http.MethodGet).Subrouter()
			getRouter.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
				id := strconv.Itoa(h.getUserId(w, r))
				w.Write(bytes.NewBufferString(id).Bytes())
			})
			getRouter.Use(h.userIdentity)

			// test
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/protected", nil)
			r.Header.Set(testCase.headerName, testCase.headerValue)

			router.ServeHTTP(w, r) // perforn request

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

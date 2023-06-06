package config_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/andredubov/todo-backend/internal/config"
)

func TestMain(t *testing.T) {

	type env struct {
		appEnv               string
		httpHost             string
		postgresHost         string
		postgresPort         string
		postgresDatabaseName string
		postgresUsername     string
		postgresPassword     string
		postgressSSLMode     string
		passwordSalt         string
		jwtSigningKey        string
	}

	type args struct {
		path string
		env  env
	}

	setEnv := func(env env) {
		os.Setenv(config.ApplicationEnvironment, env.appEnv)
		os.Setenv(config.HttpHost, env.httpHost)
		os.Setenv(config.PostgresHost, env.postgresHost)
		os.Setenv(config.PostgresPort, env.postgresPort)
		os.Setenv(config.PostgresDatabaseName, env.postgresDatabaseName)
		os.Setenv(config.PostgresUsername, env.postgresUsername)
		os.Setenv(config.PostgresPassword, env.postgresPassword)
		os.Setenv(config.PostgresSSLMode, env.postgressSSLMode)
		os.Setenv(config.PasswordSalt, env.passwordSalt)
		os.Setenv(config.JwtSigningKey, env.jwtSigningKey)
	}

	tests := []struct {
		name    string
		args    args
		want    config.Config
		wantErr bool
	}{
		{
			name: "test config",
			args: args{
				path: "../../configs",
				env: env{
					appEnv:               "local",
					httpHost:             "localhost",
					postgresHost:         "localhost",
					postgresPort:         "5432",
					postgresDatabaseName: "postgres",
					postgresUsername:     "postgres",
					postgresPassword:     "qwerty",
					postgressSSLMode:     "false",
					passwordSalt:         "salt",
					jwtSigningKey:        "key",
				},
			},
			want: config.Config{
				Environment: "local",
				CacheTTL:    time.Second * 3600,
				Postgres: config.PostgresConfig{
					Host:         "localhost",
					Port:         5432,
					Username:     "postgres",
					Password:     "qwerty",
					DatabaseName: "postgres",
					SSLMode:      "false",
				},
				HTTP: config.HTTPConfig{
					Host:               "localhost",
					MaxHeaderMegabytes: 1,
					Port:               "8080",
					ReadTimeout:        time.Second * 10,
					WriteTimeout:       time.Second * 10,
				},
				Auth: config.AuthConfig{
					PasswordSalt: "salt",
					JWT: config.JWTConfig{
						RefreshTokenTTL: time.Minute * 30,
						AccessTokenTTL:  time.Minute * 15,
						SigningKey:      "key",
					},
					VerificationCodeLength: 10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			setEnv(tt.args.env)

			got, err := config.Init(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// got  {local {localhost 5432 postgres qwerty postgres false} {{15m0s 30m0s key} salt 0} {localhost 8080 10s 10s 0} 1h0m0s}
// want {local {localhost 5432 postgres qwerty postgres false} {{15m0s 30m0s key} salt 10} {localhost 8080 10s 10s 1} 1h0m0s}

// got  {local {localhost 5432 postgres qwerty postgres false} {{0s 0s key} salt 10} {localhost 8080 10s 10s 0} 1h0m0s}
// want {local {localhost 5432 postgres qwerty postgres false} {{15m0s 30m0s key} salt 10} {localhost 8080 10s 10s 1} 1h0m0s}

// got  {local {localhost 5432 postgres qwerty postgres false} {{15m0s 30m0s key} salt 10} {localhost 8080 10s 10s 0} 1h0m0s}
// want {local {localhost 5432 postgres qwerty postgres false} {{15m0s 30m0s key} salt 10} {localhost 8080 10s 10s 1} 1h0m0s}

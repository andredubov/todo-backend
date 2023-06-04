package config

import (
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8080"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultVerificationCodeLength = 8

	Local = "local"
	Prod  = "prod"
)

type (
	Config struct {
		Environment string
		Postgres    PostgresConfig
		Auth        AuthConfig
		HTTP        HTTPConfig
		CacheTTL    time.Duration `mapstructure:"ttl"`
	}

	PostgresConfig struct {
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"port"`
		Username     string `mapstructure:"username"`
		Password     string `mapstructure:"password"`
		DatabaseName string `mapstructure:"databasename"`
		SSLMode      string `mapstructure:"sslmode"`
	}

	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeLength int `mapstructure:"verificationCodeLength"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init(configsDir string) (Config, error) {

	populateDefaults()

	var cfg Config

	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return cfg, err
	}

	if err := unmarshal(&cfg); err != nil {
		return cfg, err
	}

	if err := setFromEnv(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func unmarshal(cfg *Config) error {

	if err := viper.UnmarshalKey("cache.ttl", &cfg.CacheTTL); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) error {

	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	port, err := strconv.Atoi(os.Getenv("POSTGRESS_PORT"))
	if err != nil {
		return err
	}
	cfg.Postgres.Port = port
	cfg.Postgres.DatabaseName = os.Getenv("POSTGRESS_DB_NAME")
	cfg.Postgres.Username = os.Getenv("POSTGRES_USERNAME")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.SSLMode = os.Getenv("POSTGRES_SSL_MODE")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")
	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.Environment = os.Getenv("APP_ENV")

	return nil
}

func parseConfigFile(path, env string) error {

	viper.AddConfigPath(path)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == Local {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderBytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.readTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("auth.verificationCodeLength", defaultVerificationCodeLength)
}

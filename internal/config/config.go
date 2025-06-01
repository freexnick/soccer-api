package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const configPath = "./configs/.env"

var (
	applicationVersion string
	gitCommit          string
)

type Configuration struct {
	LogFormat string `env:"LOG_FORMAT" envDefault:"PLAIN"`
	LogLevel  string `env:"LOG_LEVEL" envDefault:"DEBUG"`

	APIVersionPrefix    string `env:"API_VERSION_PREFIX" envDefault:"/v1"`
	HTTPServerAddress   string `env:"HTTP_SERVER_ADDRESS,required"`
	ReadTimeoutSeconds  uint   `env:"HTTP_SERVER_READ_TIME_OUT" envDefault:"15"`
	WriteTimeoutSeconds uint   `env:"HTTP_SERVER_WRITE_TIME_OUT" envDefault:"45"`
	IdleTimeoutSeconds  uint   `env:"HTTP_SERVER_IDLE_TIME_OUT" envDefault:"60"`

	PostgresURL                string `env:"DB_URL,required"`
	PostgresMaxOpenConns       int    `env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"25"`
	PostgresMaxIdleConns       int    `env:"POSTGRES_MAX_IDLE_CONNS" envDefault:"10"`
	PostgresConnMaxLifetimeMin int    `env:"POSTGRES_CONN_MAX_LIFETIME_MIN" envDefault:"60"`
	PostgresConnMaxIdleTimeMin int    `env:"POSTGRES_CONN_MAX_IDLE_TIME_MIN" envDefault:"5"`

	AppVersion string `json:"version"`
	GitCommit  string `json:"git_commit"`

	JWTSecret        string `env:"JWT_SECRET,required"`
	JWTExpiryMinutes int    `env:"JWT_EXPIRY_MINUTES,required"`
	JWTIssuer        string `env:"JWT_ISSUER,required"`
}

func New() (*Configuration, error) {
	conf := &Configuration{}
	if err := godotenv.Load(configPath); err != nil {
		return nil, err
	}

	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	if applicationVersion == "" {
		conf.AppVersion = "dev"
	} else {
		conf.AppVersion = applicationVersion
	}

	if gitCommit == "" {
		conf.GitCommit = "unknown"
	} else {
		conf.GitCommit = gitCommit
	}

	return conf, nil
}

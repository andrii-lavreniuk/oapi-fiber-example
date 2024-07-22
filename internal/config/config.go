package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config interface {
	GetServiceName() string
	GetServer() Server
	GetLog() Log
	GetData() Data
}

type Server struct {
	Port        int           `env:"SERVER_PORT" envDefault:"8080"`
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"5s"`
	Debug       bool          `env:"SERVER_DEBUG" envDefault:"false"`
	OpenAPI     bool          `env:"SERVER_OPENAPI" envDefault:"false"`
}

type Log struct {
	Level  string `env:"LOG_LEVEL" envDefault:"warn"`
	Format string `env:"LOG_FORMAT" envDefault:"json"`
}

type Data struct {
	MySQL struct {
		URL     string `env:"DATA_MYSQL_URL,required"`
		Verbose bool   `env:"DATA_MYSQL_VERBOSE" envDefault:"false"`
	}
}

type config struct {
	ServiceName string `env:"SERVER_SERVICE_NAME" envDefault:"oapi-fiber-example"`
	Server      Server
	Log         Log
	Data        Data
}

func New() (Config, error) {
	cfg := &config{}

	opts := env.Options{
		Prefix: "SERVICE_",
	}

	// Load env vars.
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *config) GetServiceName() string {
	return c.ServiceName
}

func (c *config) GetServer() Server {
	return c.Server
}

func (c *config) GetLog() Log {
	return c.Log
}

func (c *config) GetData() Data {
	return c.Data
}

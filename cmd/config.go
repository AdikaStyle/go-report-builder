package cmd

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	TemplatesPath string        `env:"TEMPLATES_PATH" envDefault:"cmd/testdata/"`
	RenderTimeout time.Duration `env:"RENDER_TIMEOUT envDefault:"3s"`
	ServerHost    string        `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort    int           `env:"SERVER_PORT" envDefault:"8080"`
}

func (config Config) BaseUrl() string {
	return fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
}

func ParseConfig() (Config, error) {
	var out Config
	err := env.Parse(&out)

	if out.RenderTimeout == 0 {
		out.RenderTimeout = 3 * time.Second
	}
	return out, err
}

package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string `env:"ENV" env-default:"local"`
	Server Server
}

type Server struct {
	AppAddress      string        `env:"APP_PORT" env-default:"5555"`
	AppReadTimeout  time.Duration `env:"APP_READ_TIMEOUT" env-default:"9s"`
	AppWriteTimeout time.Duration `env:"APP_WRITE_TIMEOUT" env-default:"9s"`
	AppIdleTimeout  time.Duration `env:"APP_IDLE_TIMEOUT" env-default:"9s"`
}

// MustLoad returns Config in case no error
// If an error occurs, the app won't run and through a panic.
func MustLoad() *Config {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic("error loading environment variables")
	// }

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	cfg.Server.AppAddress = ":" + cfg.Server.AppAddress

	return &cfg
}

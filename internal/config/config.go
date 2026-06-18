package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Host         string `env:"HOST" env-default:"localhost"`
	Port         int    `env:"PORT" env-default:"8080"`
	TemplatePath string `env:"TEMPLATE_PATH" env-default:"web/index.html"`
	DBPath       string `env:"DB_PATH" env-default:"robo-ruka.db"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	_ = godotenv.Load(".env")
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

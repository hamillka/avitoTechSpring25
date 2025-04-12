package config

import (
	"github.com/hamillka/avitoTechSpring25/internal/db"
	"github.com/hamillka/avitoTechSpring25/internal/logger"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB       db.DatabaseConfig `envconfig:"DB"`
	HttpPort string            `envconfig:"HTTP_PORT"`
	GRPCPort string            `envconfig:"GRPC_PORT"`
	Timeout  int64             `envconfig:"TIMEOUT"`
	Log      logger.LogConfig  `envconfig:"LOG"`
}

func New() (*Config, error) {
	var config Config

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

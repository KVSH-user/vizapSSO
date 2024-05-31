package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Env             string         `yaml:"env" env-default:"local"`
	Postgres        PostgresConfig `yaml:"postgres"`
	AccessTokenTTL  time.Duration  `yaml:"access_token_ttl" env-required:"true"`
	RefreshTokenTTL time.Duration  `yaml:"refresh_token_ttl" env-required:"true"`
	GRPC            GRPCConfig     `yaml:"grpc"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("config path is empty")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config " + err.Error())
	}

	return &cfg
}

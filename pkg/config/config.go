package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/codingconcepts/env"
	"github.com/joho/godotenv"
)

type Config struct {
	DB       DB     `env:"DB"`
	Server   Server `env:"SERVER"`
	LogLevel string `env:"LOG_LEVEL"`
}

type Server struct {
	Host string `env:"SERVER_HOST"`
	Port string `env:"SERVER_PORT"`
}

type DB struct {
	Host string `env:"DB_HOST"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
	Name string `env:"DB_NAME"`
	Port int    `env:"DB_PORT"`
}

func Load() *Config {
	envPath := filepath.Join("..", ".", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal(err)
	}

	var config Config
	var dbConfig DB
	var serverConfig Server

	if err := env.Set(&dbConfig); err != nil {
		log.Fatal("cannot get db env vars: ", err)
	}
	if err := env.Set(&serverConfig); err != nil {
		log.Fatal("cannot get server env vars: ", err)
	}
	config.LogLevel = os.Getenv("LOG_LEVEL")

	config.DB = dbConfig
	config.Server = serverConfig

	return &config
}

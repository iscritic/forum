package config

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"forum/pkg/env"
)

type HTTPServerConfig struct {
	Address     string        `json:"address"`
	Timeout     time.Duration `json:"timeout"`
	IdleTimeout time.Duration `json:"idle_timeout"`
}

type Config struct {
	Env         string           `json:"env"`
	StoragePath string           `json:"storage_path"`
	HTTPServer  HTTPServerConfig `json:"http_server"`
}

func MustLoad() *Config {
	if err := env.LoadEnv(".env"); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not exist: %s", configPath)
	}

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	return &cfg
}

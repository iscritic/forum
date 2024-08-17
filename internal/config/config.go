package config

import (
	"encoding/json"
	"forum/pkg/flog"
	"log"
	"os"
	"time"

	"forum/pkg/env"
)

type Config struct {
	Env         string        `json:"env"`
	LogLevel    flog.LogLevel `json:"log_level"`
	StoragePath string        `json:"storage_path"`
	HTTPServer  HTTPServer    `json:"http_server"`
}

type HTTPServer struct {
	Address     string `json:"address"`
	ReadTimeout string `json:"timeout"`
	IdleTimeout string `json:"idle_timeout"`
}

func MustLoad() *Config {
	if err := env.LoadEnv(".env"); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("CONFIG_PATH is not set")
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	return &cfg
}

func ParseTime(s string) time.Duration {
	parsedTime, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("can't parse time: %v", err)
	}
	return parsedTime
}

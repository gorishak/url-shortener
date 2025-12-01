package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists and is a file
	info, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("config file does not exist: %s", configPath)
		}
		log.Fatalf("failed to stat config file: %v", err)
	}
	if info.IsDir() {
		log.Fatalf("config path is a directory, not a file: %s", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return cfg
}

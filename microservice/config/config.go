package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Mailgun  MailgunConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
}

type MailgunConfig struct {
	Domain string
	APIKey string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Database: DatabaseConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
		},
		Mailgun: MailgunConfig{
			Domain: os.Getenv("MAILGUN_DOMAIN"),
			APIKey: os.Getenv("MAILGUN_API_KEY"),
		},
	}

	cfgErr := Validate(cfg)
	if cfgErr != nil {
		log.Fatalf("Configuration error: %v", cfgErr)
	}

	return cfg, nil
}

func Validate(cfg *Config) error {
	switch {
	case cfg.Server.Port == "":
		return fmt.Errorf("PORT is required")
	case cfg.Database.User == "":
		return fmt.Errorf("DB_USER is required")
	case cfg.Database.Password == "":
		return fmt.Errorf("DB_PASSWORD is required")
	case cfg.Database.Name == "":
		return fmt.Errorf("DB_NAME is required")
	case cfg.Database.Host == "":
		return fmt.Errorf("DB_HOST is required")
	}

	return nil
}

package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server       ServerConfig
	AuthService  AuthServiceConfig
	EmailService EmailServiceConfig
	SocialAuth   SocialAuthConfig
	Database     DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type AuthServiceConfig struct {
	OTPIssuer      string
	JWTSecret      string
	AccessTokenTTL time.Duration
}

type EmailServiceConfig struct {
	Domain    string
	APIKey    string
	EmailFrom string
}

type SocialAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GithubClientID     string
	GithubClientSecret string
	GithubRedirectURL  string
}

type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
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
		AuthService: AuthServiceConfig{
			OTPIssuer:      os.Getenv("OTP_ISSUER"),
			JWTSecret:      os.Getenv("JWT_SECRET"),
			AccessTokenTTL: 15 * time.Minute,
		},
		EmailService: EmailServiceConfig{
			Domain:    os.Getenv("MAILGUN_DOMAIN"),
			APIKey:    os.Getenv("MAILGUN_API_KEY"),
			EmailFrom: os.Getenv("EMAIL_FROM"),
		},
		SocialAuth: SocialAuthConfig{
			GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			GithubRedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		},
		Database: DatabaseConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
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

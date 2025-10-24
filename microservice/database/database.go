package database

import (
	"fmt"

	"github.com/MartinLupa/secure-auth-service/microservice/config"
	"github.com/go-pg/pg/v10"
)

func InitDB(cfg *config.DatabaseConfig) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Name,
		Addr:     cfg.Host,
	})

	_, err := db.Exec("SELECT 1")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

package main

import (
	"github.com/MartinLupa/secure-auth-service/microservice/config"
	"github.com/MartinLupa/secure-auth-service/microservice/database"
	"github.com/MartinLupa/secure-auth-service/microservice/internal/handlers"
	"github.com/MartinLupa/secure-auth-service/microservice/internal/repository"
	"github.com/MartinLupa/secure-auth-service/microservice/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id       int64
	FullName string
	Email    string
	Password string
}

func main() {
	// Load configuration
	cfg, err := config.Load()

	if err != nil {
		panic("[Error] failed to load configuration due to: " + err.Error())
	}

	// Database connection
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		panic("[Error] " + err.Error())
	}
	defer db.Close()

	// Layer initialization
	authRepo := repository.NewUserRepository(db)
	emailService := service.NewEmailService(&cfg.Mailgun)
	authService := service.NewAuthService(authRepo, emailService)
	authHandler := handlers.NewAuthHandler(authService)

	// Gin router setup
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	router.POST("/login", authHandler.Login)
	router.POST("/signup", authHandler.Signup)

	err = router.Run(cfg.Server.Port)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}

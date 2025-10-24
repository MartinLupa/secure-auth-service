package handlers

import (
	"github.com/MartinLupa/secure-auth-service/microservice/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&loginData)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Email and password are required.",
		})
		return
	}

	user, err := h.authService.Login(loginData.Email, loginData.Password)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "login successful",
		"data":    user,
	})
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var signupData struct {
		FullName        string `json:"full_name" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	}

	err := c.ShouldBindJSON(&signupData)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "All fields are required and must be valid.",
		})
		return
	}

	user, err := h.authService.Signup(signupData.FullName, signupData.Email, signupData.Password, signupData.ConfirmPassword)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "signup successful",
		"data":    user,
	})
}

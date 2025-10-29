package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/MartinLupa/secure-auth-service/microservice/config"
	"github.com/MartinLupa/secure-auth-service/microservice/internal/service"
	"github.com/MartinLupa/secure-auth-service/microservice/utils"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"google.golang.org/api/idtoken"
)

type AuthHandler struct {
	config      *config.SocialAuthConfig
	authService service.AuthService
}

func NewAuthHandler(cfg *config.SocialAuthConfig, authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		config:      cfg,
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

func (h *AuthHandler) ValidateOTP(c *gin.Context) {
	var otpData struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required,len=6"`
	}

	err := c.ShouldBindJSON(&otpData)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Email and OTP are required.",
		})
		return
	}

	jwtToken, err := h.authService.ValidateOTP(otpData.Email, otpData.OTP)

	if err != nil || jwtToken == "" {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "OTP verification successful",
		"token":   jwtToken,
	})
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var requestData struct {
		Email string `json:"email" binding:"required,email"`
	}

	err := c.ShouldBindJSON(&requestData)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Valid email is required.",
		})
		return
	}

	err = h.authService.ResendOTP(requestData.Email)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "OTP resent successfully",
	})
}

func (h *AuthHandler) ValidateJWT(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(400, gin.H{
			"error": "Authorization header is required.",
		})
		return
	}

	user, err := h.authService.ValidateJWT(token)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "JWT validation successful",
		"user":    user,
	})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	goth.UseProviders(
		google.New(h.config.GoogleClientID, h.config.GoogleClientSecret, h.config.GoogleRedirectURL, "email", "profile"),
	)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to authenticate with Google: " + err.Error(),
		})
		return
	}

	payload, err := idtoken.Validate(context.Background(), gothUser.IDToken, "")

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	redirectUrl, socialUserNotFoundRedirectUrl, err := utils.ExtractDataFromUrlState(c)
	if err != nil {
		return
	}

	jwtToken, err := h.authService.GenerateSocialLoginSession(payload.Claims["email"].(string), payload.Claims["name"].(string))

	if err != nil {
		if errors.Is(err, service.ErrSocialUserNotFound) {
			c.Redirect(http.StatusFound, socialUserNotFoundRedirectUrl)
		}
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	cookie := &http.Cookie{
		Name:  "redirect_session_token",
		Value: jwtToken,
	}
	c.SetCookieData(cookie)
	c.Redirect(http.StatusFound, redirectUrl)
}

func (h *AuthHandler) GithubLogin(c *gin.Context) {
	goth.UseProviders(
		github.New(h.config.GithubClientID, h.config.GithubClientSecret, h.config.GithubRedirectURL, "email", "profile"),
	)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *AuthHandler) GithubCallback(c *gin.Context) {
	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to authenticate with Github: " + err.Error(),
		})
		return
	}

	redirectUrl, socialUserNotFoundRedirectUrl, err := utils.ExtractDataFromUrlState(c)
	if err != nil {
		return
	}

	jwtToken, err := h.authService.GenerateSocialLoginSession(gothUser.Email, gothUser.Name)

	if err != nil {
		if errors.Is(err, service.ErrSocialUserNotFound) {
			c.Redirect(http.StatusFound, socialUserNotFoundRedirectUrl)
		}
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	cookie := &http.Cookie{
		Name:  "redirect_session_token",
		Value: jwtToken,
	}
	c.SetCookieData(cookie)
	c.Redirect(http.StatusFound, redirectUrl)
}

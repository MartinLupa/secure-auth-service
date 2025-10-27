package jwt

import (
	"errors"
	"time"

	"github.com/MartinLupa/secure-auth-service/microservice/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

func GenerateAccessToken(user *models.User, jwtSecret string, accessTokenTTL time.Duration) (string, error) {
	expirationTime := time.Now().Add(accessTokenTTL)

	claims := jwt.MapClaims{
		"email": user.Email,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateAccessToken(tokenString, jwtSecret string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jwtSecret), nil
	})

	// fmt.Println("Parsed Token:", token)
	// fmt.Println("Parse Error:", err)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, ErrInvalidToken
}

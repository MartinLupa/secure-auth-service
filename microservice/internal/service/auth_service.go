package service

import (
	"errors"
	"fmt"

	"github.com/MartinLupa/secure-auth-service/microservice/config"
	"github.com/MartinLupa/secure-auth-service/microservice/internal/models"
	"github.com/MartinLupa/secure-auth-service/microservice/internal/repository"
	"github.com/MartinLupa/secure-auth-service/microservice/pkg/jwt"

	"github.com/MartinLupa/secure-auth-service/microservice/pkg/otp"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("email and password do not match")
	ErrPasswordMismatch   = errors.New("passwords do not match")
	ErrInvalidOTP         = errors.New("invalid OTP code")
	ErrGeneratingJWT      = errors.New("error generating JWT token")
	ErrSocialUserNotFound = errors.New("social login user not found")
)

type AuthService interface {
	Login(email, password string) (*models.User, error)
	Signup(fullName, email, password, confirmPassword string) (*models.User, error)
	ValidateOTP(email, code string) (string, error)
	ResendOTP(email string) error
	ValidateJWT(tokenString string) (*models.User, error)
	GenerateSocialLoginSession(email, fullName string) (string, error)
}

type authService struct {
	config       *config.AuthServiceConfig
	userRepo     repository.UserRepository
	emailService EmailService
}

func NewAuthService(cfg *config.AuthServiceConfig, userRepo repository.UserRepository, emailService EmailService) AuthService {
	return &authService{
		config:       cfg,
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *authService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	otpCode, secret, err := otp.GenerateOTP(email, s.config.OTPIssuer)

	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateOTPSecret(email, secret)

	if err != nil {
		return nil, err
	}

	fmt.Println("OTP Code: ", otpCode)
	s.emailService.SendOTPEmail(user.Email, otpCode)

	return user, nil
}

func (s *authService) Signup(fullName, email, password, confirmPassword string) (*models.User, error) {
	if password != confirmPassword {
		return nil, ErrPasswordMismatch
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	user, err := s.userRepo.CreateUser(newUser)

	if err != nil {
		return nil, err
	}

	if user != nil {
		s.emailService.SendWelcomeEmail(user.Email, fullName)
	}

	return user, nil
}

func (s *authService) ValidateOTP(email, passcode string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)

	if err != nil {
		return "", err
	}

	if user.OTPSecret == "" {
		return "", errors.New("no OTP secret found for user")
	}

	_, err = otp.ValidateOTP(passcode, user.OTPSecret)

	if err != nil {
		return "", ErrInvalidOTP
	}

	tokenString, err := jwt.GenerateAccessToken(user, s.config.JWTSecret, s.config.AccessTokenTTL)

	if err != nil {
		return "", ErrGeneratingJWT
	}

	return tokenString, nil
}

func (s *authService) ResendOTP(email string) error {
	user, err := s.userRepo.GetUserByEmail(email)

	if err != nil {
		return err
	}

	otpCode, secret, err := otp.GenerateOTP(email, s.config.OTPIssuer)

	if err != nil {
		return err
	}
	err = s.userRepo.UpdateOTPSecret(email, secret)

	if err != nil {
		return err
	}

	fmt.Println("Resent OTP Code: ", otpCode)
	s.emailService.SendOTPEmail(user.Email, otpCode)

	return nil
}

func (s *authService) ValidateJWT(tokenString string) (*models.User, error) {
	claims, err := jwt.ValidateAccessToken(tokenString, s.config.JWTSecret)

	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByEmail((*claims)["email"].(string))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) GenerateSocialLoginSession(email, fullName string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrSocialUserNotFound
		}
		return "", err
	}

	tokenString, err := jwt.GenerateAccessToken(user, s.config.JWTSecret, s.config.AccessTokenTTL)

	if err != nil {
		return "", ErrGeneratingJWT
	}

	return tokenString, nil
}

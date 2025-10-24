package service

import (
	"errors"

	"github.com/MartinLupa/secure-auth-service/microservice/internal/models"
	"github.com/MartinLupa/secure-auth-service/microservice/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("email and password do not match")
	ErrPasswordMismatch   = errors.New("passwords do not match")
)

type AuthService interface {
	Login(email, password string) (*models.User, error)
	Signup(fullName, email, password, confirmPassword string) (*models.User, error)
}

type authService struct {
	userRepo     repository.UserRepository
	emailService EmailService
}

func NewAuthService(userRepo repository.UserRepository, emailService EmailService) AuthService {
	return &authService{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *authService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	// OTP verification could be added here
	s.emailService.SendOTPEmail(user.Email, "123456") // Placeholder OTP

	return user, nil
}

func (s *authService) Signup(fullName, email, password, confirmPassword string) (*models.User, error) {
	if password != confirmPassword {
		return nil, ErrPasswordMismatch
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		FullName: fullName,
		Email:    email,
		Password: string(hashedPassword),
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

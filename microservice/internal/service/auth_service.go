package service

import (
	"errors"
	"fmt"

	"github.com/MartinLupa/secure-auth-service/microservice/internal/models"
	"github.com/MartinLupa/secure-auth-service/microservice/internal/repository"
	"github.com/MartinLupa/secure-auth-service/microservice/pkg/otp"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("email and password do not match")
	ErrPasswordMismatch   = errors.New("passwords do not match")
	ErrInvalidOTP         = errors.New("invalid OTP code")
)

type AuthService interface {
	Login(email, password string) (*models.User, error)
	Signup(fullName, email, password, confirmPassword string) (*models.User, error)
	VerifyOTP(email, code string) (bool, error)
	ResendOTP(email string) error
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

	otpCode, secret, err := otp.GenerateOTP(email, "mlupgropdevprojects@gmail.com")

	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateOTPSecret(email, secret)

	if err != nil {
		return nil, err
	}

	fmt.Println("OTP Code: ", otpCode)
	// s.emailService.SendOTPEmail(user.Email, otpCode)

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
		// s.emailService.SendWelcomeEmail(user.Email, fullName)
	}

	return user, nil
}

func (s *authService) VerifyOTP(email, code string) (bool, error) {
	user, err := s.userRepo.GetUserByEmail(email)

	if err != nil {
		return false, err
	}

	if user.OTPSecret == "" {
		return false, errors.New("no OTP secret found for user")
	}

	valid, err := otp.ValidateOTP(code, user.OTPSecret)

	if err != nil {
		return false, ErrInvalidOTP
	}

	return valid, nil
}

func (s *authService) ResendOTP(email string) error {
	// user, err := s.userRepo.GetUserByEmail(email)

	// if err != nil {
	// 	return err
	// }

	// otpCode, secret, err := otp.GenerateOTP(email, "mlupgropdevprojects@gmail.com")

	// if err != nil {
	// 	return err
	// }
	// err = s.userRepo.UpdateOTPSecret(email, secret)

	// if err != nil {
	// 	return err
	// }

	// fmt.Println("Resent OTP Code: ", otpCode)
	// s.emailService.SendOTPEmail(user.Email, otpCode)

	return nil
}

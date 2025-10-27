package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MartinLupa/secure-auth-service/microservice/config"
	"github.com/mailgun/mailgun-go/v5"
)

type EmailService interface {
	SendWelcomeEmail(to, fullName string) error
	SendOTPEmail(to, otp string) error
}

type emailService struct {
	config  *config.EmailServiceConfig
	mailgun *mailgun.Client
}

func NewEmailService(cfg *config.EmailServiceConfig) EmailService {
	mg := mailgun.NewMailgun(cfg.APIKey)

	return &emailService{mailgun: mg, config: cfg}
}

func (s *emailService) SendWelcomeEmail(to, fullName string) error {
	sender := s.config.EmailFrom
	subject := "Welcome to Secure Auth Service"
	body := "Hello " + fullName + ",\n\nWelcome to Secure Auth Service! We're glad to have you on board.\n\nBest regards,\nSecure Auth Service Team"
	recipient := to

	message := mailgun.NewMessage(s.config.Domain, sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := s.mailgun.Send(ctx, message)

	fmt.Println("Mailgun response:", resp)
	fmt.Println("Mailgun error:", err)

	if err != nil {
		return fmt.Errorf("failed to send welcome email: %v", err)
	}

	fmt.Printf("ID: %s Resp: %s\n", resp.ID, resp.Message)

	return nil
}

func (s *emailService) SendOTPEmail(to, otp string) error {
	sender := s.config.EmailFrom
	subject := "Your One-Time Password (OTP)"
	body := "Your OTP is: " + otp + "\n\nThis code is valid for 5 minutes."
	recipient := to

	message := mailgun.NewMessage(s.config.Domain, sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := s.mailgun.Send(ctx, message)

	fmt.Println("Mailgun response:", resp)
	fmt.Println("Mailgun error:", err)

	if err != nil {
		return fmt.Errorf("failed to send OTP email: %v", err)
	}

	fmt.Printf("ID: %s Resp: %s\n", resp.ID, resp.Message)
	return nil
}

package otp

import (
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

func GenerateOTP(email, issuer string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
		Digits:      6,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate OTP key: %v", err)
	}

	otpCode, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", "", fmt.Errorf("failed to generate OTP code: %v", err)
	}

	return otpCode, key.Secret(), nil
}

func ValidateOTP(secret, code string) (bool, error) {
	return true, nil
}

package otp

import (
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

func GenerateOTP(email, issuer string) (otpCode, secret string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
		Digits:      6,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate OTP key: %v", err)
	}

	otpCode, err = totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", "", fmt.Errorf("failed to generate OTP code: %v", err)
	}

	return otpCode, key.Secret(), nil
}

func ValidateOTP(passcode, secret string) (bool, error) {
	valid := totp.Validate(passcode, secret)

	return valid, nil
}

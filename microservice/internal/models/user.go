package models

type User struct {
	ID           int64  `pg:"id,pk" json:"id"`
	FullName     string `pg:"full_name" json:"full_name"`
	Email        string `pg:"email,unique" json:"email"`
	PasswordHash string `pg:"password_hash" json:"-"`
	OTPSecret    string `pg:"otp_secret" json:"-"`
}

package models

type User struct {
	ID       int64  `pg:"id,pk" json:"id"`
	FullName string `pg:"full_name" json:"full_name"`
	Email    string `pg:"email,unique" json:"email"`
	Password string `pg:"password" json:"-"`
}

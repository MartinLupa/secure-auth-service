package repository

import (
	"errors"

	"github.com/MartinLupa/secure-auth-service/microservice/internal/models"
	"github.com/go-pg/pg/v10"
)

var (
	ErrUserCreation      = errors.New("failed to create user")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	_, err := r.db.Model(user).
		Insert()

	if err != nil {
		if pgErr, ok := err.(pg.Error); ok && pgErr.IntegrityViolation() {
			return nil, ErrUserAlreadyExists
		}
		return nil, ErrUserCreation
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	err := r.db.Model(user).
		Column("full_name", "email", "password").
		Where("email = ?", email).
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, ErrUserCreation
	}

	return user, nil
}

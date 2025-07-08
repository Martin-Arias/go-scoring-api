package repository

import (
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "username = ?", username).Error
	return &user, err
}

func (r *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) RegisterUser(user *model.User) error {
	err := r.db.Create(user).Error
	return err
}

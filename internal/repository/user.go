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
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) RegisterUser(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

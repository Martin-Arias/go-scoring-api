package repository

import (
	"context"

	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) RegisterUser(user *model.User) error {
	err := r.db.Create(user).Error
	return err
}
func (r *userRepository) CreatePlayerWithInitialScores(ctx context.Context, username string, passwordHash string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		newUser := &User{
			Username:     username,
			PasswordHash: passwordHash,
		}

		if err := tx.Create(newUser).Error; err != nil {
			return err
		}

		var games []model.Game
		if err := tx.Find(&games).Error; err != nil {
			return err
		}

		var scores []model.Score
		for _, game := range games {
			scores = append(scores, model.Score{
				GameID:   game.ID,
				PlayerID: newUser.ID,
				Points:   0,
			})
		}

		if len(scores) > 0 {
			if err := tx.Create(&scores).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

package repository

import (
	"context"
	"errors"

	"github.com/Martin-Arias/go-scoring-api/internal/core/auth"
	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user User
	err := r.db.First(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &domain.User{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}, nil
}

func (r *userRepository) GetUserCreds(username string) (*auth.AuthUserData, error) {
	var user User
	err := r.db.First(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &auth.AuthUserData{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		IsAdmin:      user.IsAdmin,
	}, nil
}

func (r *userRepository) GetUserByID(id string) (*domain.User, error) {
	var user User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &domain.User{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}, nil
}

func (r *userRepository) CreateUserWithInitialScores(ctx context.Context, username string, passwordHash string) (*domain.User, error) {
	newUser := &User{
		Username:     username,
		PasswordHash: passwordHash,
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(newUser).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return domain.ErrUsernameAlreadyExists
			}
			return err
		}

		var games []Game
		if err := tx.Find(&games).Error; err != nil {
			return err
		}

		var scores []Score
		for _, game := range games {
			scores = append(scores, Score{
				GameID: game.ID,
				UserID: newUser.ID,
				Points: 0,
			})
		}

		if len(scores) > 0 {
			if err := tx.Create(&scores).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       newUser.ID,
		Username: newUser.Username,
	}, nil
}

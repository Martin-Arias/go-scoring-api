package ports

import (
	"context"

	"github.com/Martin-Arias/go-scoring-api/internal/model"
)

type UserRepository interface {
	RegisterUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	CreatePlayerWithInitialScores(ctx context.Context, username string, passwordHash string) error
}

type UserService interface {
	RegisterUser(username, password string) error
	LoginUser(username, password string) (string, error)
}

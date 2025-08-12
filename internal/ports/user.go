package ports

import (
	"context"

	"github.com/Martin-Arias/go-scoring-api/internal/core/auth"
	"github.com/Martin-Arias/go-scoring-api/internal/domain"
)

type UserRepository interface {
	GetUserByID(id string) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	GetUserCreds(username string) (*auth.AuthUserData, error)
	CreateUserWithInitialScores(ctx context.Context, username string, passwordHash string) (*domain.User, error)
}

type UserService interface {
	RegisterUser(username, password string) (*domain.User, error)
	LoginUser(username, password string) (string, error)
}

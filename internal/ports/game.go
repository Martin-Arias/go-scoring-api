package ports

import (
	"github.com/Martin-Arias/go-scoring-api/internal/domain"
)

type GameService interface {
	CreateGame(gameName string) (*domain.Game, error)
	GetGames() (*[]domain.Game, error)
}

type GameRepository interface {
	CreateGame(name string) (*domain.Game, error)
	ListGames() (*[]domain.Game, error)
	GetGameByID(id string) (*domain.Game, error)
	GetGameByName(name string) (*domain.Game, error)
}

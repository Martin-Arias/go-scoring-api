package services

import (
	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"github.com/rs/zerolog/log"
)

type gameService struct {
	gr ports.GameRepository
}

func NewGameService(gr ports.GameRepository) ports.GameService {
	return &gameService{
		gr: gr,
	}
}

func (gs *gameService) CreateGame(gameName string) (*domain.Game, error) {
	game, err := gs.gr.CreateGame(gameName)
	if err != nil {
		log.Error().Err(err).Str("game_name", gameName).Msg("failed to create game")
		return nil, err
	}

	return game, nil
}

func (gs *gameService) GetGames() (*[]domain.Game, error) {
	games, err := gs.gr.ListGames()
	if err != nil {
		log.Error().Err(err).Msg("failed retrieving games")
		return nil, err
	}

	return games, nil
}

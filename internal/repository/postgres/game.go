package repository

import (
	"context"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"gorm.io/gorm"
)

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) ports.GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) CreateGame(name string) (*domain.Game, error) {
	game := Game{Name: name}
	err := r.db.Create(&game).Error
	if err != nil {
		return nil, err
	}

	return &domain.Game{
		ID:   game.ID,
		Name: game.Name,
	}, nil
}

func (r *gameRepository) ListGames() (*[]domain.Game, error) {
	var games []Game
	err := r.db.Find(&games).Error
	if err != nil {
		return nil, err
	}

	var gamesResponse []domain.Game
	for _, game := range games {
		gamesResponse = append(gamesResponse, domain.Game{
			ID:   game.ID,
			Name: game.Name,
		})
	}
	return &gamesResponse, nil
}

func (r *gameRepository) GetGameByID(id string) (*domain.Game, error) {
	var game Game
	err := r.db.First(&game, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &domain.Game{
		ID:   game.ID,
		Name: game.Name,
	}, nil
}

func (r *gameRepository) GetGameByName(name string) (*domain.Game, error) {
	var game Game
	err := r.db.Where("name = ?", name).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &domain.Game{
		ID:   game.ID,
		Name: game.Name,
	}, nil
}

func (r *gameRepository) CreateGameWithInitialScores(ctx context.Context, name string) (*domain.Game, error) {
	newGame := &Game{
		Name: name,
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(newGame).Error; err != nil {
			return err
		}

		var users []User
		if err := tx.Find(&users).Error; err != nil {
			return err
		}

		var scores []Score
		for _, user := range users {
			scores = append(scores, Score{
				GameID:   newGame.ID,
				PlayerID: user.ID,
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

	if err != nil {
		return nil, err
	}

	return &domain.Game{
		ID:   newGame.ID,
		Name: newGame.Name,
	}, nil
}

package repository

import (
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"gorm.io/gorm"
)

type GameRepository interface {
	CreateGame(name string) (uint, error)
	ListGames() ([]model.Game, error)
	GetGameByID(id uint) (*model.Game, error)
}

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) CreateGame(name string) (uint, error) {
	game := model.Game{Name: name}
	if err := r.db.Create(&game).Error; err != nil {
		return 0, err
	}
	return game.ID, nil
}
func (r *gameRepository) ListGames() ([]model.Game, error) {
	var games []model.Game
	if err := r.db.Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}

func (r *gameRepository) GetGameByID(id uint) (*model.Game, error) {
	var game model.Game
	if err := r.db.First(&game, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &game, nil
}

package repository

import (
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"gorm.io/gorm"
)

type GameRepository interface {
	CreateGame(name string) (*model.Game, error)
	ListGames() (*[]model.Game, error)
	GetGameByID(id uint) (*model.Game, error)
	GetGameByName(name string) (*model.Game, error)
}

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) CreateGame(name string) (*model.Game, error) {
	game := model.Game{Name: name}
	err := r.db.Create(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *gameRepository) ListGames() (*[]model.Game, error) {
	var games []model.Game
	err := r.db.Find(&games).Error
	if err != nil {
		return nil, err
	}
	return &games, nil
}

func (r *gameRepository) GetGameByID(id uint) (*model.Game, error) {
	var game model.Game
	err := r.db.First(&game, id).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *gameRepository) GetGameByName(name string) (*model.Game, error) {
	var game model.Game
	err := r.db.Where("name = ?", name).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

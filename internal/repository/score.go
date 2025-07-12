package repository

import (
	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"gorm.io/gorm"
)

type ScoreRepository interface {
	GetScoresByGameID(gameID uint) (*[]dto.PlayerScoreDTO, error)
	GetScoresByPlayerID(playerID uint) (*[]dto.PlayerScoreDTO, error)
	GetScore(playerID, gameID uint) (*model.Score, error)
	SubmitScore(score model.Score) error
}

type scoreRepository struct {
	db *gorm.DB
}

func NewScoreRepository(db *gorm.DB) ScoreRepository {
	return &scoreRepository{db: db}
}

func (r *scoreRepository) GetScore(playerID, gameID uint) (*model.Score, error) {
	var score model.Score
	err := r.db.Where("player_id = ? AND game_id = ?", playerID, gameID).First(&score).Error
	if err != nil {
		return nil, err
	}
	return &score, nil
}

func (r *scoreRepository) SubmitScore(score model.Score) error {
	err := r.db.Save(&score).Error
	return err
}

func (r *scoreRepository) GetScoresByGameID(gameID uint) (*[]dto.PlayerScoreDTO, error) {
	var results []dto.PlayerScoreDTO
	err := r.db.
		Table("users").
		Select("users.username, games.name as game_name, COALESCE(scores.points, 0) as points").
		Joins("CROSS JOIN games").
		Joins("LEFT JOIN scores ON scores.player_id = users.id AND scores.game_id = games.id").
		Where("games.id = ? AND users.is_admin = false", gameID).
		Order("points DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func (r *scoreRepository) GetScoresByPlayerID(playerID uint) (*[]dto.PlayerScoreDTO, error) {
	var results []dto.PlayerScoreDTO
	err := r.db.
		Table("games").
		Select("users.username, games.name as game_name, COALESCE(scores.points, 0) as points").
		Joins("JOIN users ON users.id = ?", playerID).
		Joins("LEFT JOIN scores ON scores.game_id = games.id AND scores.player_id = users.id").
		Where("users.is_admin = false").
		Order("points DESC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}

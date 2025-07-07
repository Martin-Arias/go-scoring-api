package repository

import (
	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"gorm.io/gorm"
)

type ScoreRepository interface {
	GetScoresByGameID(gameID uint) ([]dto.PlayerScoreDTO, error)
	GetScoresByPlayerID(playerID uint) ([]dto.PlayerScoreDTO, error)
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
	if err := r.db.Where("player_id = ? AND game_id = ?", playerID, gameID).First(&score).Error; err != nil {
		return nil, err
	}
	return &score, nil
}
func (r *scoreRepository) SubmitScore(score model.Score) error {
	if err := r.db.Save(&score).Error; err != nil {
		return err
	}
	return nil
}
func (r *scoreRepository) GetScoresByGameID(gameID uint) ([]dto.PlayerScoreDTO, error) {

	var results []dto.PlayerScoreDTO
	if err := r.db.
		Table("scores").
		Select("users.username, games.name as game_name, scores.points").
		Joins("JOIN users ON users.id = scores.player_id").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("scores.game_id = ?", gameID).
		Order("scores.points DESC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *scoreRepository) GetScoresByPlayerID(playerID uint) ([]dto.PlayerScoreDTO, error) {

	var results []dto.PlayerScoreDTO
	if err := r.db.
		Table("scores").
		Select("users.username, games.name as game_name, scores.points").
		Joins("JOIN users ON users.id = scores.player_id").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("scores.player_id = ?", playerID).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

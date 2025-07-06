package repository

import (
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"gorm.io/gorm"
)

type ScoreRepository interface {
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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
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

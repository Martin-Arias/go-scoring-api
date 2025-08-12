package repository

import (
	"errors"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"gorm.io/gorm"
)

type scoreRepository struct {
	db *gorm.DB
}

func NewScoreRepository(db *gorm.DB) ports.ScoreRepository {
	return &scoreRepository{db: db}
}

func (r *scoreRepository) GetScore(userID, gameID string) (*domain.Score, error) {
	var score Score
	err := r.db.Where("user_id = ? AND game_id = ?", userID, gameID).First(&score).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrScoreNotFound
		}
		return nil, err
	}

	return &domain.Score{
		UserID: userID,
		GameID: gameID,
		Points: score.Points,
	}, nil
}

func (r *scoreRepository) SubmitScore(score *domain.Score) error {
	err := r.db.Save(&Score{
		GameID: score.GameID,
		UserID: score.UserID,
		Points: score.Points,
	}).Error
	return err
}

func (r *scoreRepository) GetScoresByGameID(gameID string) (*[]domain.Score, error) {
	var scores []dto.UserScoreDTO
	err := r.db.
		Table("scores").
		Select("users.username, scores.user_id, games.name as game_name, scores.game_id, scores.points").
		Joins("JOIN users ON users.id = scores.user_id").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("scores.game_id = ?", gameID).
		Order("scores.points DESC").
		Scan(&scores).Error
	if err != nil {
		return nil, err
	}

	if len(scores) == 0 {
		return nil, domain.ErrScoreNotFound
	}

	var scoresResponse []domain.Score
	for _, score := range scores {
		scoresResponse = append(scoresResponse, domain.Score{
			Username: score.Username,
			UserID:   score.UserID,
			GameName: score.GameName,
			GameID:   score.GameID,
			Points:   score.Points,
		})
	}

	return &scoresResponse, nil
}

func (r *scoreRepository) GetScoresByUserID(playerID string) (*[]domain.Score, error) {
	var scores []dto.UserScoreDTO
	err := r.db.
		Table("scores").
		Select("users.username, scores.user_id, games.name as game_name, scores.game_id, scores.points").
		Joins("JOIN users ON users.id = scores.user_id").
		Joins("JOIN games ON games.id = scores.game_id").
		Order("scores.points DESC").
		Where("scores.user_id = ?", playerID).
		Scan(&scores).Error
	if err != nil {
		return nil, err
	}

	if len(scores) == 0 {
		return nil, domain.ErrScoreNotFound
	}

	var scoresResponse []domain.Score
	for _, score := range scores {
		scoresResponse = append(scoresResponse, domain.Score{
			Username: score.Username,
			UserID:   score.UserID,
			GameName: score.GameName,
			GameID:   score.GameID,
			Points:   score.Points,
		})
	}

	return &scoresResponse, nil
}

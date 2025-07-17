package repository

import (
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

func (r *scoreRepository) GetScore(playerID, gameID string) (*domain.Score, error) {
	var score Score
	err := r.db.Where("player_id = ? AND game_id = ?", playerID, gameID).First(&score).Error
	if err != nil {
		return nil, err
	}

	return &domain.Score{
		ID:     score.ID,
		Points: score.Points,
	}, nil
}

func (r *scoreRepository) SubmitScore(score *domain.Score) error {
	err := r.db.Save(&score).Error
	return err
}

func (r *scoreRepository) GetScoresByGameID(gameID string) (*[]domain.Score, error) {
	var scores []dto.PlayerScoreDTO //Move this dto to score_model
	err := r.db.
		Table("scores").
		Select("users.username, games.name as game_name, scores.points").
		Joins("JOIN users ON users.id = scores.player_id").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("scores.game_id = ?", gameID).
		Order("scores.points DESC").
		Scan(&scores).Error
	if err != nil {
		return nil, err
	}

	var scoresResponse []domain.Score
	for _, score := range scores {
		scoresResponse = append(scoresResponse, domain.Score{
			User: domain.User{
				Username: score.Username,
			},
			Game: domain.Game{
				Name: score.GameName,
			},
			Points: score.Points,
		})
	}
	return &scoresResponse, nil
}

func (r *scoreRepository) GetScoresByPlayerID(playerID string) (*[]domain.Score, error) {
	var scores []dto.PlayerScoreDTO
	err := r.db.
		Table("scores").
		Select("users.username, games.name as game_name, scores.points").
		Joins("JOIN users ON users.id = scores.player_id").
		Joins("JOIN games ON games.id = scores.game_id").
		Where("scores.player_id = ?", playerID).
		Scan(&scores).Error
	if err != nil {
		return nil, err
	}
	var scoresResponse []domain.Score
	for _, score := range scores {
		scoresResponse = append(scoresResponse, domain.Score{
			User: domain.User{
				Username: score.Username,
			},
			Game: domain.Game{
				Name: score.GameName,
			},
			Points: score.Points,
		})
	}
	return &scoresResponse, nil
}

package ports

import (
	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/Martin-Arias/go-scoring-api/internal/dto"
)

type ScoreRepository interface {
	GetScoresByGameID(gameID string) (*[]domain.Score, error)
	GetScoresByUserID(playerID string) (*[]domain.Score, error)
	GetScore(playerID, gameID string) (*domain.Score, error)
	SubmitScore(score *domain.Score) error
}

type ScoreService interface {
	Submit(score *domain.Score) error
	GetGameScores(gameID string) (*[]domain.Score, error)
	GetUserScores(userID string) (*[]domain.Score, error)
	GetGameStats(gameID string) (*dto.ScoreStatisticsDTO, error)
}

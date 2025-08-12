package mocks

import (
	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type ScoreRepositoryMock struct {
	mock.Mock
}

func (m *ScoreRepositoryMock) GetScore(userID, gameID string) (*domain.Score, error) {
	args := m.Called(userID, gameID)
	return args.Get(0).(*domain.Score), args.Error(1)
}
func (m *ScoreRepositoryMock) SubmitScore(score *domain.Score) error {
	args := m.Called(score)
	return args.Error(0)
}
func (m *ScoreRepositoryMock) GetScoresByGameID(gameID string) (*[]domain.Score, error) {
	args := m.Called(gameID)
	return args.Get(0).(*[]domain.Score), args.Error(1)
}
func (m *ScoreRepositoryMock) GetScoresByUserID(userID string) (*[]domain.Score, error) {
	args := m.Called(userID)
	return args.Get(0).(*[]domain.Score), args.Error(1)
}

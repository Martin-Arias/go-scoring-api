package mocks

/*
import (
	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type ScoreRepositoryMock struct {
	mock.Mock
}

func (m *ScoreRepositoryMock) GetScoresByGameID(gameID string) (*[]dto.PlayerScoreDTO, error) {
	args := m.Called(gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]dto.PlayerScoreDTO), args.Error(1)
}

func (m *ScoreRepositoryMock) GetScoresByPlayerID(playerID string) (*[]dto.PlayerScoreDTO, error) {
	args := m.Called(playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]dto.PlayerScoreDTO), args.Error(1)
}

func (m *ScoreRepositoryMock) GetScore(playerID, gameID string) (*model.Score, error) {
	args := m.Called(playerID, gameID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Score), args.Error(1)
}

func (m *ScoreRepositoryMock) SubmitScore(score model.Score) error {
	args := m.Called(score)
	return args.Error(0)
}
*/

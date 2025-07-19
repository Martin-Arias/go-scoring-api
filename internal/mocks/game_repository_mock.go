package mocks

import (
	"context"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type GameRepositoryMock struct {
	mock.Mock
}

func (m *GameRepositoryMock) CreateGameWithInitialScores(ctx context.Context, name string) (*domain.Game, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*domain.Game), args.Error(1)
}

func (m *GameRepositoryMock) ListGames() (*[]domain.Game, error) {
	args := m.Called()
	return args.Get(0).(*[]domain.Game), args.Error(1)
}

func (m *GameRepositoryMock) GetGameByID(id string) (*domain.Game, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Game), args.Error(1)
}

func (m *GameRepositoryMock) GetGameByName(name string) (*domain.Game, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Game), args.Error(1)
}

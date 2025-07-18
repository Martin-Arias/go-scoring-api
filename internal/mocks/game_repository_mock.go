package mocks

/*
import (
	"github.com/stretchr/testify/mock"
)

type GameRepositoryMock struct {
	mock.Mock
}

func (m *GameRepositoryMock) ListGames() (*[]model.Game, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]model.Game), args.Error(1)
}

func (m *GameRepositoryMock) CreateGame(name string) (*model.Game, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}

func (m *GameRepositoryMock) GetGameByID(id string) (*model.Game, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}

func (m *GameRepositoryMock) GetGameByName(name string) (*model.Game, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}
*/

package mocks

import (
	"context"

	"github.com/Martin-Arias/go-scoring-api/internal/core/auth"
	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetUserByID(userID string) (*domain.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByUsername(userID string) (*domain.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUserCreds(username string) (*auth.AuthUserData, error) {
	args := m.Called(username)
	return args.Get(0).(*auth.AuthUserData), args.Error(1)
}

func (m *UserRepositoryMock) CreateUserWithInitialScores(ctx context.Context, username string, passwordHash string) (*domain.User, error) {
	args := m.Called(ctx, username, passwordHash)
	return args.Get(0).(*domain.User), args.Error(1)
}

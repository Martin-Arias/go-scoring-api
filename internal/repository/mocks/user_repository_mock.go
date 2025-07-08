package mocks

import (
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetUserByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByID(id uint) (*model.User, error) {
	args := m.Called(id)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}
func (m *UserRepositoryMock) RegisterUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

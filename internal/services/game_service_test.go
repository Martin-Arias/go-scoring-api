package services_test

import (
	"testing"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	mocks "github.com/Martin-Arias/go-scoring-api/internal/mocks/repository"
	"github.com/Martin-Arias/go-scoring-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGame(t *testing.T) {
	mockRepo := new(mocks.GameRepositoryMock)
	service := services.NewGameService(mockRepo)

	expected := &domain.Game{ID: "123", Name: "chess"}
	mockRepo.On("CreateGameWithInitialScores", mock.Anything, "chess").Return(expected, nil)

	game, err := service.CreateGame("chess")
	assert.NoError(t, err)
	assert.Equal(t, expected, game)

	mockRepo.AssertExpectations(t)
}

func TestCreateGame_Error(t *testing.T) {
	mockRepo := new(mocks.GameRepositoryMock)
	service := services.NewGameService(mockRepo)

	var nilGame *domain.Game = nil

	mockRepo.
		On("CreateGameWithInitialScores", mock.Anything, "fail-game").
		Return(nilGame, assert.AnError)

	game, err := service.CreateGame("fail-game")

	assert.Error(t, err)
	assert.Nil(t, game)
	mockRepo.AssertExpectations(t)
}

func TestGetGames(t *testing.T) {
	mockRepo := new(mocks.GameRepositoryMock)
	service := services.NewGameService(mockRepo)

	games := []domain.Game{{ID: "1", Name: "checkers"}, {ID: "2", Name: "pong"}}
	mockRepo.On("ListGames").Return(&games, nil)

	result, err := service.GetGames()
	assert.NoError(t, err)
	assert.Equal(t, &games, result)

	mockRepo.AssertExpectations(t)
}

func TestGetGames_Error(t *testing.T) {
	mockRepo := new(mocks.GameRepositoryMock)
	service := services.NewGameService(mockRepo)

	var nilGames *[]domain.Game = nil

	mockRepo.
		On("ListGames").
		Return(nilGames, assert.AnError)

	games, err := service.GetGames()

	assert.Error(t, err)
	assert.Nil(t, games)
	mockRepo.AssertExpectations(t)
}

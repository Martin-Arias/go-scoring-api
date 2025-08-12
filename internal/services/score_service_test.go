package services_test

import (
	"errors"
	"testing"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	mocks "github.com/Martin-Arias/go-scoring-api/internal/mocks/repository"
	"github.com/Martin-Arias/go-scoring-api/internal/services"
	"github.com/stretchr/testify/assert"
)

var validUser = &domain.User{ID: "user1", Username: "test", IsAdmin: false}
var validGame = &domain.Game{ID: "game1", Name: "testgame"}

var newScore = &domain.Score{UserID: "user1", GameID: "game1", Points: 101}
var validScore = &domain.Score{UserID: "user1", GameID: "game1", Points: 100}

func TestSubmitScore_NewScoreSuccess(t *testing.T) {
	sr := new(mocks.ScoreRepositoryMock)
	ur := new(mocks.UserRepositoryMock)
	gr := new(mocks.GameRepositoryMock)

	ss := services.NewScoreService(sr, ur, gr)

	ur.On("GetUserByID", "user1").Return(validUser, nil)
	gr.On("GetGameByID", "game1").Return(validGame, nil)
	sr.On("GetScore", "user1", "game1").Return(validScore, nil)
	sr.On("SubmitScore", newScore).Return(nil)

	err := ss.Submit(newScore)
	assert.NoError(t, err)

	sr.AssertExpectations(t)
	ur.AssertExpectations(t)
	gr.AssertExpectations(t)
}

func TestSubmitScore_UserNotFound(t *testing.T) {
	sr := new(mocks.ScoreRepositoryMock)
	ur := new(mocks.UserRepositoryMock)
	gr := new(mocks.GameRepositoryMock)

	ss := services.NewScoreService(sr, ur, gr)

	ur.On("GetUserByID", "user1").Return(&domain.User{}, errors.New("user not found"))

	err := ss.Submit(validScore)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestSubmitScore_GameNotFound(t *testing.T) {
	sr := new(mocks.ScoreRepositoryMock)
	ur := new(mocks.UserRepositoryMock)
	gr := new(mocks.GameRepositoryMock)

	ss := services.NewScoreService(sr, ur, gr)

	ur.On("GetUserByID", "user1").Return(validUser, nil)
	gr.On("GetGameByID", "game1").Return(nil, errors.New("game not found"))

	err := ss.Submit(validScore)
	assert.Error(t, err)
	assert.Equal(t, "game not found", err.Error())
}

func TestSubmitScore_LowerScoreRejected(t *testing.T) {
	sr := new(mocks.ScoreRepositoryMock)
	ur := new(mocks.UserRepositoryMock)
	gr := new(mocks.GameRepositoryMock)

	ss := services.NewScoreService(sr, ur, gr)

	oldScore := &domain.Score{UserID: "user1", GameID: "game1", Points: 200}

	ur.On("GetUserByID", "user1").Return(validUser, nil)
	gr.On("GetGameByID", "game1").Return(validGame, nil)
	sr.On("GetScore", "user1", "game1").Return(oldScore, nil)

	err := ss.Submit(validScore)
	assert.ErrorIs(t, err, domain.ErrScoreNotAllowed)
}

func TestGetUserScores_Success(t *testing.T) {
	sr := new(mocks.ScoreRepositoryMock)
	ur := new(mocks.UserRepositoryMock)
	gr := new(mocks.GameRepositoryMock)

	ss := services.NewScoreService(sr, ur, gr)

	userScores := &[]domain.Score{
		{GameID: "game1", Points: 100},
	}

	ur.On("GetUserByID", "user1").Return(validUser, nil)
	sr.On("GetScoresByUserID", "user1").Return(userScores, nil)

	scores, err := ss.GetUserScores("user1")
	assert.NoError(t, err)
	assert.Equal(t, userScores, scores)
}

func TestGetGameStats_Success(t *testing.T) {
	sr := new(mocks.ScoreRepositoryMock)
	ur := new(mocks.UserRepositoryMock)
	gr := new(mocks.GameRepositoryMock)

	ss := services.NewScoreService(sr, ur, gr)

	scoreList := &[]domain.Score{
		{GameID: "game1", GameName: "Test Game", Points: 10},
		{GameID: "game1", GameName: "Test Game", Points: 20},
		{GameID: "game1", GameName: "Test Game", Points: 10},
	}

	sr.On("GetScoresByGameID", "game1").Return(scoreList, nil)

	stats, err := ss.GetGameStats("game1")
	assert.NoError(t, err)
	assert.Equal(t, "game1", stats.GameID)
	assert.Equal(t, "Test Game", stats.GameName)
	assert.Equal(t, 13.33, stats.Mean)
}

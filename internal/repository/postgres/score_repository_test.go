package repository_test

import (
	"context"
	"testing"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	repository "github.com/Martin-Arias/go-scoring-api/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
)

func TestScoreRepository_CreateAndQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := repository.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)
	scoreRepo := repository.NewScoreRepository(db)

	game, err := gameRepo.CreateGameWithInitialScores(context.Background(), "pong")
	assert.NoError(t, err)
	t.Logf("Created game: ID=%s, Name=%s", game.ID, game.Name)
	// setup
	user, err := userRepo.CreateUserWithInitialScores(context.Background(), "juan", "123")
	assert.NoError(t, err)
	t.Logf("Created user: ID=%s, Username=%s", user.ID, user.Username)

	err = scoreRepo.SubmitScore(&domain.Score{
		GameID: game.ID,
		UserID: user.ID,
		Points: 1000,
	})
	assert.NoError(t, err)

	scores, err := scoreRepo.GetScoresByGameID(game.ID)
	assert.NoError(t, err)
	assert.Len(t, *scores, 1)

	for _, s := range *scores {
		t.Logf("Score: UserID=%s, GameID=%s, Value=%d", s.UserID, s.GameID, s.Points)
	}
}

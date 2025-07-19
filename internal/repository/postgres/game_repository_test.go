package repository_test

import (
	"context"
	"testing"

	repository "github.com/Martin-Arias/go-scoring-api/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
)

func TestGameRepository_CreateAndFetch(t *testing.T) {
	db := repository.SetupTestDB(t)
	repo := repository.NewGameRepository(db)

	game, err := repo.CreateGameWithInitialScores(context.Background(), "game-1")
	assert.NoError(t, err)
	assert.NotNil(t, game)

	byID, err := repo.GetGameByID(game.ID)
	assert.NoError(t, err)
	assert.Equal(t, game.ID, byID.ID)

	byName, err := repo.GetGameByName(game.Name)
	assert.NoError(t, err)
	assert.Equal(t, game.Name, byName.Name)

	allGames, err := repo.ListGames()
	assert.NoError(t, err)
	assert.Len(t, *allGames, 1)
}

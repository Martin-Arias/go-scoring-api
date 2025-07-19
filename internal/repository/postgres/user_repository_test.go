package repository_test

import (
	"context"
	"testing"

	repository "github.com/Martin-Arias/go-scoring-api/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateAndFetch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := repository.SetupTestDB(t)
	repo := repository.NewUserRepository(db)

	_, err := repo.CreateUserWithInitialScores(context.Background(), "martin", "pass123")
	assert.NoError(t, err)

	user, err := repo.GetUserByUsername("martin")
	assert.NoError(t, err)
	assert.Equal(t, "martin", user.Username)
}

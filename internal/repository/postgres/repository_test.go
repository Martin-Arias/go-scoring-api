package repository_test

import (
	"context"
	"testing"
	"time"

	repository "github.com/Martin-Arias/go-scoring-api/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()
	t.Helper()
	// Start PostgreSQL container
	container, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := gorm.Open(gormPostgres.Open(connStr), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	repository.RunMigrations(db)

	return db
}

func TestCreateAndGetGame(t *testing.T) {
	db := setupTestDB(t)

	repo := repository.NewGameRepository(db)

	// Create a game
	game, err := repo.CreateGame("test-game")
	assert.NoError(t, err)
	assert.NotNil(t, game)
	assert.Equal(t, "test-game", game.Name)

	// Get the game
	fetched, err := repo.GetGameByID(game.ID)
	assert.NoError(t, err)
	assert.Equal(t, game.ID, fetched.ID)
	assert.Equal(t, "test-game", fetched.Name)
}

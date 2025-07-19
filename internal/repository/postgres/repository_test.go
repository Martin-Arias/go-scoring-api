package repository_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	repository "github.com/Martin-Arias/go-scoring-api/internal/repository/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpass",
			"POSTGRES_DB=testdb",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	dsn := fmt.Sprintf("host=localhost port=%s user=testuser password=testpass dbname=testdb sslmode=disable", resource.GetPort("5432/tcp"))

	if err := pool.Retry(func() error {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
		if err != nil {
			return err
		}

		db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
		// Automatically migrate the schema
		if err := db.AutoMigrate(&repository.User{}, &repository.Score{}, &repository.Game{}); err != nil {
			return fmt.Errorf("error running migration: %w", err)
		}
		db.Exec(`ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4();`)
		db.Exec(`ALTER TABLE games ALTER COLUMN id SET DEFAULT uuid_generate_v4();`)

		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestCreateGame(t *testing.T) {
	r := repository.NewGameRepository(db)

	game, err := r.CreateGame("Test Game")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if game.Name != "Test Game" {
		t.Errorf("expected game name 'Test Game', got '%s'", game.Name)
	}
}

func TestCreateGame_Duplicated(t *testing.T) {
	r := repository.NewGameRepository(db)

	_, _ = r.CreateGame("Duplicate Game")
	_, err := r.CreateGame("Duplicate Game")
	if err != domain.ErrGameAlreadyExists {
		t.Errorf("expected ErrGameAlreadyExists, got %v", err)
	}
}

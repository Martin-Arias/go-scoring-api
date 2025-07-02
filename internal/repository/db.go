package repository

import (
	"fmt"

	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndMigrate(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema
	if err := db.AutoMigrate(&model.User{}, &model.Score{}, &model.Game{}); err != nil {
		return nil, fmt.Errorf("error migrando: %w", err)
	}

	return db, nil
}

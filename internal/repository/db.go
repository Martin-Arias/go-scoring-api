package repository

import (
	"fmt"
	"os"

	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndMigrate() (*gorm.DB, error) {
	//dsn
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema
	if err := db.AutoMigrate(&model.User{}, &model.Score{}, &model.Game{}); err != nil {
		return nil, fmt.Errorf("error running migration: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error creating admin user: %w", err)
	}

	// Create the admin user if it doesn't exist
	var adminUser model.User
	if err := db.First(&adminUser, "username = ?", "admin").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			adminUser = model.User{
				Username:     "admin",
				PasswordHash: string(hash), // hashed password for "admin123"
				IsAdmin:      true,
			}
			if err := db.Create(&adminUser).Error; err != nil {
				return nil, fmt.Errorf("error creating admin user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error checking for admin user: %w", err)
		}
	}

	return db, nil
}

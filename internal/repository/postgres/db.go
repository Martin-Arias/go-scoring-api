package repository

import (
	"fmt"
	"os"

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

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Automatically migrate the schema
	if err := db.AutoMigrate(&User{}, &Score{}, &Game{}); err != nil {
		return nil, fmt.Errorf("error running migration: %w", err)
	}
	db.Exec(`ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4();`)
	db.Exec(`ALTER TABLE games ALTER COLUMN id SET DEFAULT uuid_generate_v4();`)

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error creating admin user: %w", err)
	}

	// Create the admin user if it doesn't exist
	var adminUser User
	if err := db.First(&adminUser, "username = ?", "admin").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			adminUser = User{
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

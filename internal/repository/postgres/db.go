package repository

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	//dsn
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Second) // No limit on connection lifetime

	RunMigrations(db)

	return db, nil
}

func RunMigrations(db *gorm.DB) error {

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return fmt.Errorf("failed to create extension: %w", err)
	}

	if err := db.AutoMigrate(&User{}, &Score{}, &Game{}); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := db.Exec(`ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4()`).Error; err != nil {
		return err
	}

	if err := db.Exec(`ALTER TABLE games ALTER COLUMN id SET DEFAULT uuid_generate_v4()`).Error; err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error creating admin user: %w", err)
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
				return fmt.Errorf("error creating admin user: %w", err)
			}
		} else {
			return fmt.Errorf("error creating admin user: %w", err)
		}
	}
	return nil
}

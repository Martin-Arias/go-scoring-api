package repository

import (
	"time"

	"github.com/Martin-Arias/go-scoring-api/internal/model"
)

type User struct {
	ID           string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	IsAdmin      bool   `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	//FK
	Scores []model.Score `gorm:"foreignKey:PlayerID;constraint:OnDelete:CASCADE"`
}

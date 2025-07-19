package repository

import (
	"time"
)

type User struct {
	ID           string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	IsAdmin      bool   `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	//FK
	Scores []Score `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

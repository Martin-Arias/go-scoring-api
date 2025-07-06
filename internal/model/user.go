package model

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"` // Para este id lo correcto seria un UUID, pero me es mas fácil copiar y pegar uint
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	IsAdmin      bool   `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	//FK
	Scores []Score `gorm:"foreignKey:PlayerID;constraint:OnDelete:CASCADE"`
}

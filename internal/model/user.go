package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	IsAdmin      bool   `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	//FK
	Scores []Score `gorm:"foreignKey:PlayerID;constraint:OnDelete:CASCADE"`
}

// AfterCreate crea puntuaciones por defecto para todos los juegos existentes
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.IsAdmin {
		return nil
	}
	var games []Game
	if err := tx.Find(&games).Error; err != nil {
		return err
	}
	for _, game := range games {
		score := Score{
			PlayerID: u.ID,
			GameID:   game.ID,
			Points:   0,
		}
		if err := tx.Create(&score).Error; err != nil {
			return err
		}
	}
	return nil
}

package model

import "gorm.io/gorm"

type Game struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`

	//FK
	Scores []Score `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
}

// AfterCreate crea puntuaciones 0 para todos los jugadores existentes
func (g *Game) AfterCreate(tx *gorm.DB) (err error) {
	var users []User
	if err := tx.Where("is_admin = ?", false).Find(&users).Error; err != nil {
		return err
	}

	for _, player := range users {
		score := Score{
			PlayerID: player.ID,
			GameID:   g.ID,
			Points:   0,
		}
		if err := tx.Create(&score).Error; err != nil {
			return err
		}
	}

	return nil
}

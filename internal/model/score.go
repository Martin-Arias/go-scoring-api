package model

type Score struct {
	ID       uint `gorm:"primaryKey"`
	PlayerID uint `gorm:"primaryKey"`
	GameID   uint `gorm:"primaryKey"`
	Points   int  `gorm:"not null"`

	// FKs
	Player User `gorm:"foreignKey:PlayerID"`
	Game   Game `gorm:"foreignKey:GameID"`
}

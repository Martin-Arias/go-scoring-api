package repository

type Score struct {
	UserID string `gorm:"primaryKey"`
	GameID string `gorm:"primaryKey"`
	Points int    `gorm:"not null"`

	// FKs
	User User `gorm:"foreignKey:UserID"`
	Game Game `gorm:"foreignKey:GameID"`
}

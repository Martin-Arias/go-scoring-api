package repository

type Score struct {
	ID     string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID string `gorm:"primaryKey"`
	GameID string `gorm:"primaryKey"`
	Points int    `gorm:"not null"`

	// FKs
	User User `gorm:"foreignKey:UserID"`
	Game Game `gorm:"foreignKey:GameID"`
}

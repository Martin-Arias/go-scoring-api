package repository

type Score struct {
	ID       string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PlayerID string `gorm:"primaryKey"`
	GameID   string `gorm:"primaryKey"`
	Points   int    `gorm:"not null"`

	// FKs
	Player User `gorm:"foreignKey:PlayerID"`
	Game   Game `gorm:"foreignKey:GameID"`
}

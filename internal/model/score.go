package model

type Score struct {
	PlayerID string `gorm:"primaryKey"`
	GameID   string `gorm:"primaryKey"`
	Points   int    `gorm:"not null"`

	// FKs
	Player User `gorm:"foreignKey:PlayerID"`
	Game   Game `gorm:"foreignKey:GameID"`
}

type Game struct {
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`

	//FK
	Scores []Score `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
}

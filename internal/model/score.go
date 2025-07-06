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

type Game struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`

	//FK
	Scores []Score `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
}

package model

import "time"

type Score struct {
	PlayerID  string    `gorm:"primaryKey"`
	GameID    string    `gorm:"primaryKey;index:idx_scores_game_points"`
	Timestamp time.Time `gorm:"primaryKey"`

	Points int `gorm:"not null;index:idx_scores_game_points,sort:desc"`

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

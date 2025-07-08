package model

type Game struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`

	//FK
	Scores []Score `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
}

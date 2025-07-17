package repository

type Game struct {
	ID   string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name string `gorm:"uniqueIndex;not null"`

	//FK
	Scores []Score `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
}

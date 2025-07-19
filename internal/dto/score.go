package dto

type UserScoreDTO struct {
	UserID   string `json:"user_id"   gorm:"column:user_id"`
	GameID   string `json:"game_id"   gorm:"column:game_id"`
	Username string `json:"username"  gorm:"column:username"`
	GameName string `json:"game_name" gorm:"column:game_name"`
	Points   int    `json:"points"    gorm:"column:points"`
}

type ScoreStatisticsDTO struct {
	GameID   string  `json:"game_id"`
	GameName string  `json:"game_name"`
	Mean     float64 `json:"mean"`
	Median   float64 `json:"median"`
	Mode     []int   `json:"mode"`
}

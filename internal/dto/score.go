package dto

type PlayerScoreDTO struct {
	Username string `json:"username"`
	GameName string `json:"game_name"`
	Points   int    `json:"points"`
}

type ScoreStatisticsDTO struct {
	GameID   string  `json:"game_id"`
	GameName string  `json:"game_name"`
	Mean     float64 `json:"mean"`
	Median   float64 `json:"median"`
	Mode     []int   `json:"mode"`
}

package dto

type PlayerScoreDTO struct {
	Username string `json:"username"`
	GameName string `json:"game_name"`
	Points   int    `json:"points"`
}

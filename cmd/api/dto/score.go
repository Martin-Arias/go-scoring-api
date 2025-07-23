package dto

type SubmitScoreRequest struct {
	UserID string `json:"user_id" binding:"required,uuid4"`
	GameID string `json:"game_id" binding:"required,uuid4"`
	Points int    `json:"points" binding:"required,min=0"`
}

type ScoreResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	GameID   string `json:"game_id"`
	GameName string `json:"game_name"`
	Points   int    `json:"points"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ScoreStatisticsDTO struct {
	GameID   string  `json:"game_id"`
	GameName string  `json:"game_name"`
	Mean     float64 `json:"mean"`
	Median   float64 `json:"median"`
	Mode     []int   `json:"mode"`
}

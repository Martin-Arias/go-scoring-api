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

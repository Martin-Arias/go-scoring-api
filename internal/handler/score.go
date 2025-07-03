package handler

import (
	"gorm.io/gorm"
)

type SubmitScoreRequest struct {
	PlayerID string `json:"player_id" binding:"required"`
	GameID   string `json:"game_id" binding:"required"`
	Points   int    `json:"points" binding:"required,gte=0"`
}

type ScoreHandler struct {
	db *gorm.DB
}

func NewScoreHandler(db *gorm.DB) *ScoreHandler {
	return &ScoreHandler{db: db}
}

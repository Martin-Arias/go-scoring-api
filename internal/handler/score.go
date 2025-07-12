package handler

import (
	"net/http"
	"strconv"

	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/Martin-Arias/go-scoring-api/internal/repository"
	"github.com/Martin-Arias/go-scoring-api/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubmitScoreRequest struct {
	PlayerID uint `json:"player_id" binding:"required"`
	GameID   uint `json:"game_id" binding:"required"`
	Points   int  `json:"points" binding:"required,gte=0"`
}

type ScoreHandler struct {
	sr repository.ScoreRepository
	ur repository.UserRepository
	gr repository.GameRepository
}

func NewScoreHandler(sr repository.ScoreRepository, ur repository.UserRepository, gr repository.GameRepository) *ScoreHandler {
	return &ScoreHandler{sr: sr, ur: ur, gr: gr}
}

func (h *ScoreHandler) Submit(c *gin.Context) {

	var req SubmitScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Check if user exists
	usr, err := h.ur.GetUserByID(req.PlayerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Admin users aren't allowed to play :'(
	if usr.IsAdmin {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}

	// Check if game exists
	_, err = h.gr.GetGameByID(req.GameID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	score := model.Score{
		PlayerID: req.PlayerID,
		GameID:   req.GameID,
		Points:   req.Points,
	}

	existingScore, err := h.sr.GetScore(req.PlayerID, req.GameID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Check if score points are less than new score points
	if existingScore != nil && existingScore.Points >= req.Points {
		c.JSON(http.StatusConflict, gin.H{"error": "score must be higher than existing score"})
		return
	}

	if existingScore != nil {
		score.ID = existingScore.ID // If score exists, use its ID to update
	}

	if err := h.sr.SubmitScore(score); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "score submitted successfully"})
}

func (h *ScoreHandler) GetScoresByGameID(c *gin.Context) {
	gameIDStr := c.Query("game_id")
	if gameIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	scores, err := h.sr.GetScoresByGameID(uint(gameID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, scores)

}

func (h *ScoreHandler) GetScoresByPlayerID(c *gin.Context) {
	playerIDStr := c.Query("player_id")
	if playerIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
		return
	}

	playerID, err := strconv.Atoi(playerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
		return
	}

	scores, err := h.sr.GetScoresByPlayerID(uint(playerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, scores)

}

func (h *ScoreHandler) GetStatisticsByGameID(c *gin.Context) {

	gameIDStr := c.Query("game_id")
	if gameIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	scores, err := h.sr.GetScoresByGameID(uint(gameID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(*scores) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no scores for this game"})
		return
	}

	points := []int{}
	for _, s := range *scores {
		points = append(points, s.Points)
	}

	mean, median, mode := utils.CalculateStatistics(points)

	c.JSON(http.StatusOK, dto.ScoreStatisticsDTO{
		GameID:   uint(gameID),
		GameName: (*scores)[0].GameName,
		Mean:     mean,
		Median:   median,
		Mode:     mode,
	})
}

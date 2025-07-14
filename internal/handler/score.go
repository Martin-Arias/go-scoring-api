package handler

import (
	"net/http"
	"strconv"

	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/Martin-Arias/go-scoring-api/internal/repository"
	"github.com/Martin-Arias/go-scoring-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type SubmitScoreRequest struct {
	PlayerID uint `json:"player_id" binding:"required"`
	GameID   uint `json:"game_id" binding:"required"`
	Points   int  `json:"points" binding:"required"`
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
		log.Warn().Err(err).Msg("invalid submit score request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	log.Debug().Uint("player_id", req.PlayerID).Uint("game_id", req.GameID).Int("points", req.Points).Msg("submitting score")

	usr, err := h.ur.GetUserByID(req.PlayerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Uint("player_id", req.PlayerID).Msg("player not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
			return
		}
		log.Error().Err(err).Msg("error fetching player")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if usr.IsAdmin {
		log.Warn().Str("username", usr.Username).Msg("admin tried to submit score")
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}

	_, err = h.gr.GetGameByID(req.GameID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Uint("game_id", req.GameID).Msg("game not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}
		log.Error().Err(err).Msg("error fetching game")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	existingScore, err := h.sr.GetScore(req.PlayerID, req.GameID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("error retrieving existing score")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if existingScore != nil && existingScore.Points >= req.Points {
		log.Info().
			Uint("player_id", req.PlayerID).
			Uint("game_id", req.GameID).
			Int("existing_points", existingScore.Points).
			Int("new_points", req.Points).
			Msg("score not updated - new score is not higher")
		c.JSON(http.StatusConflict, gin.H{"error": "score must be higher than existing score"})
		return
	}

	score := model.Score{
		PlayerID: req.PlayerID,
		GameID:   req.GameID,
		Points:   req.Points,
	}
	if existingScore != nil {
		score.ID = existingScore.ID
	}

	if err := h.sr.SubmitScore(score); err != nil {
		log.Error().Err(err).Msg("failed to submit score")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit score"})
		return
	}

	log.Info().Uint("player_id", req.PlayerID).Uint("game_id", req.GameID).Int("points", req.Points).Msg("score submitted successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "score submitted successfully"})
}

func (h *ScoreHandler) GetScoresByGameID(c *gin.Context) {
	gameIDStr := c.Query("game_id")
	if gameIDStr == "" {
		log.Warn().Msg("missing game_id in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		log.Warn().Str("game_id", gameIDStr).Msg("invalid game ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	_, err = h.gr.GetGameByID(uint(gameID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Uint("game_id", uint(gameID)).Msg("game not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}
		log.Error().Err(err).Msg("error checking game existence")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	scores, err := h.sr.GetScoresByGameID(uint(gameID))
	if err != nil {
		log.Error().Err(err).Uint("game_id", uint(gameID)).Msg("error retrieving scores by game")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Info().Uint("game_id", uint(gameID)).Int("count", len(*scores)).Msg("scores retrieved successfully")
	c.JSON(http.StatusOK, scores)
}

func (h *ScoreHandler) GetScoresByPlayerID(c *gin.Context) {
	playerIDStr := c.Query("player_id")
	if playerIDStr == "" {
		log.Warn().Msg("missing player_id in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
		return
	}

	playerID, err := strconv.Atoi(playerIDStr)
	if err != nil {
		log.Warn().Str("player_id", playerIDStr).Msg("invalid player ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
		return
	}

	_, err = h.ur.GetUserByID(uint(playerID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Uint("player_id", uint(playerID)).Msg("player not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
			return
		}
		log.Error().Err(err).Msg("error checking user existence")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	scores, err := h.sr.GetScoresByPlayerID(uint(playerID))
	if err != nil {
		log.Error().Err(err).Uint("player_id", uint(playerID)).Msg("error retrieving scores by player")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Info().Uint("player_id", uint(playerID)).Int("count", len(*scores)).Msg("scores retrieved successfully")
	c.JSON(http.StatusOK, scores)
}

func (h *ScoreHandler) GetStatisticsByGameID(c *gin.Context) {
	gameIDStr := c.Query("game_id")
	if gameIDStr == "" {
		log.Warn().Msg("missing game_id in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		log.Warn().Str("game_id", gameIDStr).Msg("invalid game ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	scores, err := h.sr.GetScoresByGameID(uint(gameID))
	if err != nil {
		log.Error().Err(err).Uint("game_id", uint(gameID)).Msg("error retrieving scores for statistics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(*scores) == 0 {
		log.Info().Uint("game_id", uint(gameID)).Msg("no scores found for statistics")
		c.JSON(http.StatusNotFound, gin.H{"error": "no scores for this game"})
		return
	}

	points := make([]int, len(*scores))
	for i, s := range *scores {
		points[i] = s.Points
	}

	mean, median, mode := utils.CalculateStatistics(points)

	log.Info().
		Uint("game_id", uint(gameID)).
		Float64("mean", mean).
		Float64("median", median).
		Ints("mode", mode).
		Msg("score statistics calculated")

	c.JSON(http.StatusOK, dto.ScoreStatisticsDTO{
		GameID:   uint(gameID),
		GameName: (*scores)[0].GameName,
		Mean:     mean,
		Median:   median,
		Mode:     mode,
	})
}

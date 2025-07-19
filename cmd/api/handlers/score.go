package handlers

import (
	"errors"
	"net/http"

	"github.com/Martin-Arias/go-scoring-api/cmd/api/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ScoreHandler struct {
	ss ports.ScoreService
}

func NewScoreHandler(ss ports.ScoreService) *ScoreHandler {
	return &ScoreHandler{ss: ss}
}

// Submit submits or updates a user's score for a game.
//
// @Summary Submit a score
// @Description Submits or updates the score for a user in a specific game
// @Tags scores
// @Accept json
// @Produce json
// @Param request body SubmitScoreRequest true "Score data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]interface{} "error: string"
// @Failure 409 {object} map[string]interface{} "error: string"
// @Failure 500 {object} map[string]interface{} "error: string"
// @Security BearerAuth
// @Router /api/scores [put]
func (h *ScoreHandler) Submit(c *gin.Context) {
	var req dto.SubmitScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("invalid submit score request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	log.Debug().Str("user_id", req.UserID).Str("game_id", req.GameID).Int("points", req.Points).Msg("submitting score")

	err := h.ss.Submit(&domain.Score{
		GameID: req.GameID,
		UserID: req.UserID,
		Points: req.Points,
	})

	if err != nil {
		log.Warn().Err(err).Any("req", req).Msg("score could not be submitted")
		switch {
		case errors.Is(err, domain.ErrScoreNotAllowed):
			c.JSON(http.StatusConflict, gin.H{"error": domain.ErrScoreNotAllowed.Error()})
			return

		case errors.Is(err, domain.ErrGameNotFound), errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed submitting score"})
			return
		}
	}

	log.Info().Str("user_id", req.UserID).Str("game_id", req.GameID).Int("points", req.Points).Msg("score submitted successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "score submitted successfully"})
}

// GetGameScores returns all scores for a given game.
//
// @Summary Get scores by game
// @Description Lists user scores for a specific game
// @Tags scores
// @Produce json
// @Param game_id query int true "Game ID"
// @Success 200 {array} dto.PlayerScoreDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/scores/game [get]
func (h *ScoreHandler) GetGameScores(c *gin.Context) {
	gameID := c.Query("game_id")
	if gameID == "" {
		log.Warn().Msg("missing game_id in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	scores, err := h.ss.GetGameScores(gameID)
	if err != nil {
		log.Warn().Err(err).Msg("game scores could not be retrieved")
		if errors.Is(err, domain.ErrGameNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrGameNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieving game scores"})
		return
	}

	var response []dto.ScoreResponse
	for _, score := range *scores {
		response = append(response, dto.ScoreResponse{
			ID:       score.ID,
			UserID:   score.UserID,
			Username: score.Username,
			GameID:   score.GameID,
			GameName: score.GameName,
			Points:   score.Points,
		})
	}

	log.Info().Str("game_id", gameID).Int("count", len(*scores)).Msg("scores retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// GetUserScores returns all scores for a specific User.
//
// @Summary Get scores by User
// @Description Lists game scores for a specific User
// @Tags scores
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} dto.PlayerScoreDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/scores/user [get]
func (h *ScoreHandler) GetUserScores(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		log.Warn().Msg("missing user_id in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	scores, err := h.ss.GetUserScores(userID)
	if err != nil {
		log.Warn().Err(err).Msg("user scores could not be retrieved")
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieving user scores"})
		return
	}

	var response []dto.ScoreResponse
	for _, score := range *scores {
		response = append(response, dto.ScoreResponse{
			ID:       score.ID,
			UserID:   score.UserID,
			Username: score.Username,
			GameID:   score.GameID,
			GameName: score.GameName,
			Points:   score.Points,
		})
	}

	log.Info().Str("user_id", userID).Int("count", len(*scores)).Msg("user scores retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// GetStatisticsByGameID returns score statistics (mean, median, mode) for a game.
//
// @Summary Get game score statistics
// @Description Calculates mean, median, and mode for a game's scores
// @Tags scores
// @Produce json
// @Param game_id query int true "Game ID"
// @Success 200 {object} dto.ScoreStatisticsDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/scores/game/stats [get]
func (h *ScoreHandler) GetGameStats(c *gin.Context) {
	gameID := c.Query("game_id")
	if gameID == "" {
		log.Warn().Msg("missing game_id in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	stats, err := h.ss.GetGameStats(gameID)
	if err != nil {
		log.Warn().Err(err).Msg("game stats could not be retrieved")
		if errors.Is(err, domain.ErrScoreNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrScoreNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed retrieving game stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

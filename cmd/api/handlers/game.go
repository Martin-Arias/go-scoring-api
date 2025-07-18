package handlers

import (
	"errors"
	"net/http"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type GameHandler struct {
	gs ports.GameService
}

func NewGameHandler(gs ports.GameService) *GameHandler {
	return &GameHandler{gs: gs}
}

// Create creates a new game.
//
// @Summary Create a game
// @Description Adds a new game to the system
// @Tags games
// @Accept json
// @Produce json
// @Param request body dto.GameDTO true "Game name"
// @Success 201 {object} dto.GameDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /api/games [post]
func (h *GameHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("invalid input for game creation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	createdGame, err := h.gs.CreateGame(req.Name)
	if err != nil {
		log.Warn().Err(err).Str("name", req.Name).Msg("game could not be created")
		if errors.Is(err, domain.ErrGameAlreadyExists) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": domain.ErrGameAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed creating game"})
		return
	}

	log.Info().Str("game_id", createdGame.ID).Str("game_name", createdGame.Name).Msg("game created successfully")
	c.JSON(http.StatusCreated, createdGame)
}

// List returns all games.
//
// @Summary List games
// @Description Retrieves all available games
// @Tags games
// @Produce json
// @Success 200 {array} dto.GameDTO
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/games [get]
func (h *GameHandler) List(c *gin.Context) {
	games, err := h.gs.GetGames()
	if err != nil {
		log.Error().Err(err).Msg("error listing games")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Info().Int("game_count", len(*games)).Msg("games listed successfully")
	c.JSON(http.StatusOK, games)
}

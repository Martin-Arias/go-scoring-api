package handler

import (
	"net/http"

	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type GameHandler struct {
	gr repository.GameRepository
}

func NewGameHandler(gr repository.GameRepository) *GameHandler {
	return &GameHandler{gr: gr}
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

	log.Debug().Str("game_name", req.Name).Msg("checking if game already exists")
	existingGame, err := h.gr.GetGameByName(req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Str("game_name", req.Name).Msg("error checking existing game")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if existingGame != nil {
		log.Info().Str("game_name", req.Name).Msg("game already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "game already exists"})
		return
	}

	game, err := h.gr.CreateGame(req.Name)
	if err != nil {
		log.Error().Err(err).Str("game_name", req.Name).Msg("failed to create game")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Info().Uint("game_id", game.ID).Str("game_name", game.Name).Msg("game created successfully")
	c.JSON(http.StatusCreated, dto.GameDTO{ID: game.ID, Name: game.Name})
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
	games, err := h.gr.ListGames()
	if err != nil {
		log.Error().Err(err).Msg("error listing games")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Info().Int("game_count", len(*games)).Msg("games listed successfully")

	var response []dto.GameDTO
	for _, game := range *games {
		response = append(response, dto.GameDTO{ID: game.ID, Name: game.Name})
	}

	c.JSON(http.StatusOK, response)
}

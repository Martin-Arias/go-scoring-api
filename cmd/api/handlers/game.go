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

type GameHandler struct {
	gs ports.GameService
}

func NewGameHandler(gs ports.GameService) *GameHandler {
	return &GameHandler{gs: gs}
}

// Create creates a new game.
//
// @Summary Create a new game
// @Description Adds a new game to the system with a unique name.
// @Tags games
// @Accept json
// @Produce json
// @Param request body dto.CreateRequest true "Game to create"
// @Success 201 {object} dto.GameResponse "Game created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 409 {object} map[string]string "Game already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/games [post]
func (h *GameHandler) Create(c *gin.Context) {

	var createReq dto.CreateRequest

	if err := c.BindJSON(&createReq); err != nil {
		log.Warn().Err(err).Msg("invalid input for game creation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	createdGame, err := h.gs.CreateGame(createReq.Name)
	if err != nil {
		log.Warn().Err(err).Str("name", createReq.Name).Msg("game could not be created")
		if errors.Is(err, domain.ErrGameAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": domain.ErrGameAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed creating game"})
		return
	}

	log.Info().Str("game_id", createdGame.ID).Str("game_name", createdGame.Name).Msg("game created successfully")
	c.JSON(http.StatusCreated, dto.GameResponse{
		ID:   createdGame.ID,
		Name: createdGame.Name,
	})
}

// List returns all games.
//
// @Summary Get list of games
// @Description Retrieves all games available in the system.
// @Tags games
// @Produce json
// @Success 200 {array} dto.GameResponse "List of games"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/games [get]
func (h *GameHandler) List(c *gin.Context) {
	games, err := h.gs.GetGames()
	if err != nil {
		log.Error().Err(err).Msg("error listing games")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var response []dto.GameResponse
	for _, game := range *games {
		response = append(response, dto.GameResponse{
			ID:   game.ID,
			Name: game.Name,
		})
	}
	log.Info().Int("game_count", len(*games)).Msg("games listed successfully")
	c.JSON(http.StatusOK, response)
}

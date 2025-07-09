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

func (h *GameHandler) Create(c *gin.Context) {

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	//check if game already exists
	existingGame, err := h.gr.GetGameByName(req.Name)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}
	if existingGame != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "game already exists"})
		return
	}

	game, err := h.gr.CreateGame(req.Name)

	if err != nil {
		log.Error().Err(err).Msg("error creating game")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, dto.GameDTO{ID: game.ID, Name: game.Name})
}

func (h *GameHandler) List(c *gin.Context) {
	games, err := h.gr.ListGames()
	if err != nil {
		log.Error().Err(err).Msg("error listing games")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var response []dto.GameDTO
	for _, game := range *games {
		response = append(response, dto.GameDTO{ID: game.ID, Name: game.Name})
	}

	c.JSON(http.StatusOK, response)
}

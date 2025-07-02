package handler

import (
	"net/http"

	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type GameHandler struct {
	db *gorm.DB
}

func NewGameHandler(db *gorm.DB) *GameHandler {
	return &GameHandler{db: db}
}

func (h *GameHandler) Create(c *gin.Context) {

	var req struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	game := model.Game{
		ID:   uuid.New().String(),
		Name: req.Name,
	}

	if err := h.db.Create(&game).Error; err != nil {
		log.Error().Err(err).Msg("error creating game")
		c.JSON(http.StatusConflict, gin.H{"error": "error creating game"})
		return
	}

	c.JSON(http.StatusCreated, dto.GameResponse{ID: game.ID, Name: game.Name})
}

func (h *GameHandler) List(c *gin.Context) {
	var games []model.Game
	if err := h.db.Find(&games).Error; err != nil {
		log.Error().Err(err).Msg("error fetching games")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching games"})
		return
	}

	var result []dto.GameResponse
	for _, g := range games {
		result = append(result, dto.GameResponse{ID: g.ID, Name: g.Name})
	}

	c.JSON(http.StatusOK, result)
}

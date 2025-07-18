package handlers

import (
	"net/http"

	"github.com/Martin-Arias/go-scoring-api/cmd/api/core"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	us ports.UserService
}

func NewUserHandler(us ports.UserService) *UserHandler {
	return &UserHandler{us: us}
}

// Register creates a new user account.
//
// @Summary Register a new user
// @Description Creates a user with a username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "User credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]interface{} "error: string"
// @Failure 409 {object} map[string]interface{} "error: string"
// @Failure 500 {object} map[string]interface{} "error: string"
// @Router /auth/register [post]
func (uh *UserHandler) Register(c *gin.Context) {
	var req core.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("invalid register request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if _, err := uh.us.RegisterUser(req.Username, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

// Login authenticates a user and returns a JWT.
//
// @Summary Login user
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "User credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	var req core.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := uh.us.LoginUser(req.Username, req.Password)
	if err != nil {
		log.Warn().Err(err).Msg("error getting user")
		c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	log.Info().Str("username", req.Username).Msg("user logged in successfully")
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

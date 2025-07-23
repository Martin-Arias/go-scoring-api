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
// @Param request body dto.AuthRequest true "User credentials"
// @Success 201 {object} dto.RegisterResponse "User registered successfully"
// @Failure 400 {object} map[string]interface{} "error: Invalid request"
// @Failure 409 {object} map[string]interface{} "error: Username already exists"
// @Failure 500 {object} map[string]interface{} "error: Internal error"
// @Router /auth/register [post]
func (uh *UserHandler) Register(c *gin.Context) {
	var req dto.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("invalid register request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	createdUser, err := uh.us.RegisterUser(req.Username, req.Password)
	if err != nil {
		log.Error().Err(err).Str("user_name", req.Username).Msg("failed to register user")
		if errors.Is(err, domain.ErrUsernameAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": domain.ErrUsernameAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error registering user"})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username,
	})
}

// Login authenticates a user and returns a JWT.
//
// @Summary Login user
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.AuthRequest true "User credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]interface{} "error: Invalid request"
// @Failure 404 {object} map[string]interface{} "error: User not found"
// @Failure 500 {object} map[string]interface{} "error: Internal error"
// @Router /auth/login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	var req dto.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := uh.us.LoginUser(req.Username, req.Password)
	if err != nil {
		log.Warn().Err(err).Str("username", req.Username).Msg("failed to login user")

		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrAuthInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrAuthInvalid.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Info().Str("username", req.Username).Msg("user logged in successfully")
	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
	})
}

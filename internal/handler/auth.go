package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/Martin-Arias/go-scoring-api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthHandler struct {
	ur repository.UserRepository
}

func NewAuthHandler(ur repository.UserRepository) *AuthHandler {
	return &AuthHandler{ur: ur}
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
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("invalid register request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.ur.GetUserByUsername(req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("error checking existing user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if user != nil {
		log.Info().Str("username", req.Username).Msg("username already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("error hashing password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	newUser := model.User{
		Username:     req.Username,
		PasswordHash: string(hash),
	}

	if err := h.ur.RegisterUser(&newUser); err != nil {
		log.Error().Err(err).Str("username", req.Username).Msg("failed to register user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	log.Info().Str("username", req.Username).Msg("user registered successfully")
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
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.ur.GetUserByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Info().Str("username", req.Username).Msg("login failed: user not found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		log.Error().Err(err).Str("username", req.Username).Msg("error fetching user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Info().Str("username", req.Username).Msg("login failed: invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := GenerateToken(*user)
	if err != nil {
		log.Error().Err(err).Str("username", req.Username).Msg("failed to generate token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Info().Str("username", user.Username).Uint("user_id", user.ID).Msg("user logged in successfully")
	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": user.Username,
	})
}

func GenerateToken(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"uid":      user.ID,
		"username": user.Username,
		"admin":    user.IsAdmin,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

package services

import (
	"context"
	"os"
	"time"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	ur ports.UserRepository
}

func NewUserService(ur ports.UserRepository) ports.UserService {
	return &UserService{
		ur: ur,
	}
}

func (us *UserService) RegisterUser(username, password string) (*domain.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return nil, err
	}

	createdUser, err := us.ur.CreatePlayerWithInitialScores(context.Background(), username, string(hash))
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("failed to register user")
		return nil, err
	}

	return createdUser, nil
}

func (us *UserService) LoginUser(username, password string) (string, error) {
	user, err := us.ur.GetUserCreds(username)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("failed to fetch user")
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Info().Str("username", username).Msg("login failed: invalid password")
		return "", domain.ErrAuthInvalid
	}

	token, err := GenerateToken(&domain.User{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	})

	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("failed to generate token")
		return "", domain.ErrUnexpected
	}

	return token, nil
}

func GenerateToken(user *domain.User) (string, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
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

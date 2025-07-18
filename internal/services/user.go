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
	"gorm.io/gorm"
)

type UserService struct {
	ur ports.UserRepository
}

func NewUserService(ur ports.UserRepository) ports.UserService {
	return &UserService{
		ur: ur,
	}
}

func (us *UserService) RegisterUser(username, password string) error {
	user, err := us.ur.GetUserByUsername(username)

	//FIX: User domain errors
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("error checking existing user")
		return err
	}

	if user != nil {
		log.Info().Str("username", username).Msg("username already exists")
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("error hashing password")
		return err
	}

	if err := us.ur.CreatePlayerWithInitialScores(context.Background(), username, string(hash)); err != nil {
		log.Error().Err(err).Str("username", username).Msg("failed to register user")
		return err
	}

	return nil
}

func (us *UserService) LoginUser(username, password string) (string, error) {
	user, err := us.ur.GetUserCreds(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Info().Str("username", username).Msg("login failed: user not found")
			return "", err
		}
		log.Error().Err(err).Str("username", username).Msg("error fetching user")
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Info().Str("username", username).Msg("login failed: invalid password")
		return "", err
	}

	token, err := GenerateToken(&domain.User{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	})

	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("failed to generate token")
		return "", nil
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

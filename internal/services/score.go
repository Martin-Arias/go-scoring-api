package services

import (
	"errors"

	"github.com/Martin-Arias/go-scoring-api/internal/domain"
	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/ports"
	"github.com/Martin-Arias/go-scoring-api/internal/utils"
	"github.com/rs/zerolog/log"
)

type ScoreService struct {
	sr ports.ScoreRepository
	ur ports.UserRepository
	gr ports.GameRepository
}

func NewScoreService(sr ports.ScoreRepository, ur ports.UserRepository, gr ports.GameRepository) ports.ScoreService {
	return &ScoreService{
		sr: sr,
		ur: ur,
		gr: gr,
	}
}

func (ss *ScoreService) Submit(newScore *domain.Score) error {

	usr, err := ss.ur.GetUserByID(newScore.UserID)
	if err != nil {
		log.Error().Err(err).Str("user_id", newScore.UserID).Msg("error fetching user")
		return err
	}

	if usr.IsAdmin {
		log.Warn().Str("username", usr.Username).Msg("admin tried to submit score")
		return domain.ErrUserNotFound
	}

	_, err = ss.gr.GetGameByID(newScore.GameID)
	if err != nil {
		log.Error().Err(err).Str("game_id", newScore.GameID).Msg("error fetching game")
		return err
	}

	existingScore, err := ss.sr.GetScore(newScore.UserID, newScore.GameID)
	if err != nil {
		if !errors.Is(err, domain.ErrScoreNotFound) {
			log.Error().Err(err).Any("newScore", newScore).Msg("error checking if score existence")
			return err
		}
	}

	if existingScore != nil && existingScore.Points >= newScore.Points {
		log.Info().
			Str("user_id", newScore.UserID).
			Str("game_id", newScore.GameID).
			Int("existing_points", existingScore.Points).
			Int("new_points", newScore.Points).
			Msg("score not updated - new score is not higher")
		return domain.ErrScoreNotAllowed
	}

	if existingScore != nil {
		newScore.ID = existingScore.ID
	}

	if err := ss.sr.SubmitScore(newScore); err != nil {
		log.Error().Err(err).Msg("failed to submit score")
		return err
	}

	return nil
}

func (ss *ScoreService) GetGameScores(gameID string) (*[]domain.Score, error) {
	_, err := ss.gr.GetGameByID(gameID)
	if err != nil {
		log.Error().Err(err).Str("game_id", gameID).Msg("error checking game existence")
		return nil, err
	}

	scores, err := ss.sr.GetScoresByGameID(gameID)
	if err != nil {
		log.Error().Err(err).Str("game_id", gameID).Msg("error retrieving scores by game")
		return nil, err
	}

	return scores, nil

}

func (ss *ScoreService) GetUserScores(userID string) (*[]domain.Score, error) {
	_, err := ss.ur.GetUserByID(userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("error fetching user")
		return nil, err
	}

	scores, err := ss.sr.GetScoresByUserID(userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("error retrieving user scores")
		return nil, err
	}

	return scores, nil
}

func (ss *ScoreService) GetGameStats(gameID string) (*dto.ScoreStatisticsDTO, error) {

	scores, err := ss.sr.GetScoresByGameID(gameID)
	if err != nil {
		log.Error().Err(err).Str("game_id", gameID).Msg("error retrieving scores for statistics")
		return nil, err
	}

	points := make([]int, len(*scores))
	for i, s := range *scores {
		points[i] = s.Points
	}

	mean, median, mode := utils.CalculateStatistics(points)

	log.Info().
		Str("game_id", gameID).
		Float64("mean", mean).
		Float64("median", median).
		Ints("mode", mode).
		Msg("score statistics calculated")

	return &dto.ScoreStatisticsDTO{
		GameID:   gameID,
		GameName: (*scores)[0].GameName,
		Mean:     mean,
		Median:   median,
		Mode:     mode,
	}, nil
}

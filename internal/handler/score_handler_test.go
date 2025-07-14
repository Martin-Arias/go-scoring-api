package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/handler"
	"github.com/Martin-Arias/go-scoring-api/internal/mocks"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var _ = Describe("ScoreHandler", func() {
	var (
		router *gin.Engine
		sr     *mocks.ScoreRepositoryMock
		ur     *mocks.UserRepositoryMock
		gr     *mocks.GameRepositoryMock
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.Default()

		sr = new(mocks.ScoreRepositoryMock)
		ur = new(mocks.UserRepositoryMock)
		gr = new(mocks.GameRepositoryMock)

		handler := handler.NewScoreHandler(sr, ur, gr)
		router.PUT("/scores", handler.Submit)
		router.GET("/scores/game", handler.GetScoresByGameID)
		router.GET("/scores/user", handler.GetScoresByPlayerID)
		router.GET("/scores/game/stats", handler.GetStatisticsByGameID)
	})

	Describe("Submit score", func() {
		It("should return 201 when score is submitted successfully", func() {
			reqBody := `{"player_id":1, "game_id":2, "points":100}`

			ur.On("GetUserByID", uint(1)).Return(&model.User{ID: 1, Username: "mario"}, nil)
			gr.On("GetGameByID", uint(2)).Return(&model.Game{ID: 2, Name: "Chess"}, nil)
			sr.On("GetScore", uint(1), uint(2)).Return(&model.Score{ID: 10, Points: 60}, nil)
			sr.On("SubmitScore", mock.Anything).Return(nil)

			req := httptest.NewRequest(http.MethodPut, "/scores", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(rec.Body.String()).To(ContainSubstring("score submitted successfully"))
		})

		It("should return 409 if submitted score is not higher", func() {
			reqBody := `{"player_id":1, "game_id":2, "points":50}`

			ur.On("GetUserByID", uint(1)).Return(&model.User{ID: 1}, nil)
			gr.On("GetGameByID", uint(2)).Return(&model.Game{ID: 2}, nil)
			sr.On("GetScore", uint(1), uint(2)).Return(&model.Score{ID: 10, Points: 60}, nil)

			req := httptest.NewRequest(http.MethodPut, "/scores", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusConflict))
			Expect(rec.Body.String()).To(ContainSubstring("score must be higher"))
		})

		It("should return 404 if player not exists", func() {
			reqBody := `{"player_id":1, "game_id":2, "points":50}`

			ur.On("GetUserByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)

			req := httptest.NewRequest(http.MethodPut, "/scores", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusNotFound))
			Expect(rec.Body.String()).To(ContainSubstring("player not found"))
		})

		It("should return 404 if game not exists", func() {
			reqBody := `{"player_id":1, "game_id":2, "points":50}`

			ur.On("GetUserByID", uint(1)).Return(&model.User{ID: 1}, nil)
			gr.On("GetGameByID", uint(2)).Return(nil, gorm.ErrRecordNotFound)

			req := httptest.NewRequest(http.MethodPut, "/scores", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusNotFound))
			Expect(rec.Body.String()).To(ContainSubstring("game not found"))
		})
	})

	Describe("GetScoresByGameID", func() {
		It("should return list of scores", func() {
			scores := []dto.PlayerScoreDTO{
				{Username: "mario", GameName: "Chess", Points: 100},
				{Username: "luigi", GameName: "Chess", Points: 80},
			}

			gr.On("GetGameByID", uint(2)).Return(&model.Game{}, nil)
			sr.On("GetScoresByGameID", uint(2)).Return(&scores, nil)

			req := httptest.NewRequest(http.MethodGet, "/scores/game?game_id=2", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("mario"))
			Expect(rec.Body.String()).To(ContainSubstring("luigi"))
		})

		It("should return 404 if game not exists", func() {
			gr.On("GetGameByID", uint(2)).Return(nil, gorm.ErrRecordNotFound)

			req := httptest.NewRequest(http.MethodGet, "/scores/game?game_id=2", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GetStatisticsByGameID", func() {
		It("should return statistics correctly", func() {
			scores := []dto.PlayerScoreDTO{
				{Username: "mario", GameName: "Chess", Points: 100},
				{Username: "luigi", GameName: "Chess", Points: 80},
				{Username: "peach", GameName: "Chess", Points: 100},
			}

			sr.On("GetScoresByGameID", uint(2)).Return(&scores, nil)

			req := httptest.NewRequest(http.MethodGet, "/scores/game/stats?game_id=2", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("mean"))
			Expect(rec.Body.String()).To(ContainSubstring("mode"))
			Expect(rec.Body.String()).To(ContainSubstring("median"))
		})

		It("should return 404 if no scores", func() {
			sr.On("GetScoresByGameID", uint(2)).Return(&[]dto.PlayerScoreDTO{}, nil)

			req := httptest.NewRequest(http.MethodGet, "/scores/game/stats?game_id=2", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GetScoresByPlayerID", func() {
		It("should return list of scores", func() {
			scores := []dto.PlayerScoreDTO{
				{Username: "mario", GameName: "Age of empires", Points: 100},
				{Username: "mario", GameName: "Chess", Points: 80},
				{Username: "mario", GameName: "Tetris", Points: 0},
			}

			ur.On("GetUserByID", uint(2)).Return(&model.User{}, nil)
			sr.On("GetScoresByPlayerID", uint(2)).Return(&scores, nil)

			req := httptest.NewRequest(http.MethodGet, "/scores/user?player_id=2", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("mario"))
		})

		It("should return 404 if user not exists", func() {
			ur.On("GetUserByID", uint(2)).Return(nil, gorm.ErrRecordNotFound)

			req := httptest.NewRequest(http.MethodGet, "/scores/user?player_id=2", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusNotFound))
		})
	})
})

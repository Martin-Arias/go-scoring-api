package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/handler"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/Martin-Arias/go-scoring-api/internal/repository/mocks"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var _ = Describe("GameHandler List", func() {
	var (
		r        *gin.Engine
		mockRepo *mocks.GameRepositoryMock
		gameH    *handler.GameHandler
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		r = gin.Default()
		mockRepo = new(mocks.GameRepositoryMock)
		gameH = handler.NewGameHandler(mockRepo)
		r.GET("/games", gameH.List)
	})

	Context("when games exist", func() {
		It("returns 200 and a list of games", func() {
			games := &[]model.Game{
				{ID: 1, Name: "Tetris"},
				{ID: 2, Name: "Age of Empires"},
			}
			mockRepo.On("ListGames").Return(games, nil)

			req, _ := http.NewRequest(http.MethodGet, "/games", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusOK))

			var response []dto.GameDTO
			err := json.Unmarshal(resp.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response).To(HaveLen(2))
			Expect(response[0].Name).To(Equal("Tetris"))
		})
	})

	Context("when repository returns an error listing games", func() {
		It("returns 500", func() {
			mockRepo.On("ListGames").Return(nil, assert.AnError)

			req, _ := http.NewRequest(http.MethodGet, "/games", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})

var _ = Describe("GameHandler Create", func() {
	var (
		r        *gin.Engine
		mockRepo *mocks.GameRepositoryMock
		gameH    *handler.GameHandler
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		r = gin.Default()
		mockRepo = new(mocks.GameRepositoryMock)
		gameH = handler.NewGameHandler(mockRepo)
		r.POST("/games", gameH.Create)
	})

	Context("when a new game is created", func() {
		It("returns 201 and the created game json", func() {
			createdGame := &model.Game{
				ID: 1, Name: "Tetris",
			}
			body := map[string]string{
				"name": "Tetris",
			}
			jsonBody, _ := json.Marshal(body)

			mockRepo.On("GetGameByName", body["name"]).Return(nil, gorm.ErrRecordNotFound)
			mockRepo.On("CreateGame", body["name"]).Return(createdGame, nil)

			req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusCreated))

			var response dto.GameDTO
			err := json.Unmarshal(resp.Body.Bytes(), &response)
			Expect(err).To(BeNil())
			Expect(response.Name).To(Equal("Tetris"))
		})
	})

	Context("when the game to create already exists", func() {
		It("returns 409 - conflict error", func() {
			existingGame := &model.Game{
				ID: 1, Name: "Tetris",
			}
			body := map[string]string{
				"name": "Tetris",
			}
			jsonBody, _ := json.Marshal(body)

			mockRepo.On("GetGameByName", body["name"]).Return(existingGame, nil)

			req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusConflict))
		})
	})

	Context("when a game creation fails", func() {
		It("returns 500", func() {
			body := map[string]string{
				"name": "Tetris",
			}
			jsonBody, _ := json.Marshal(body)

			mockRepo.On("GetGameByName", body["name"]).Return(nil, gorm.ErrRecordNotFound)
			mockRepo.On("CreateGame", body["name"]).Return(nil, errors.New("some error"))

			req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusInternalServerError))

		})
	})
})

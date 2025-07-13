package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Martin-Arias/go-scoring-api/internal/dto"
	"github.com/Martin-Arias/go-scoring-api/internal/handler"
	"github.com/Martin-Arias/go-scoring-api/internal/mocks"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("GameHandler", func() {
	var (
		router *gin.Engine
		gr     *mocks.GameRepositoryMock
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.Default()
		gr = new(mocks.GameRepositoryMock)
		handler := handler.NewGameHandler(gr)

		router.GET("/games", handler.List)
		router.POST("/games", handler.Create)
	})

	Describe("GET /games", func() {
		It("should return 200 and a list of games", func() {
			gr.On("ListGames").Return(&[]model.Game{
				{ID: 1, Name: "Tetris"},
				{ID: 2, Name: "Age of Empires"},
			}, nil)

			resp := makeRequest(router, http.MethodGet, "/games", nil)

			Expect(resp.Code).To(Equal(http.StatusOK))

			var body []dto.GameDTO
			err := json.Unmarshal(resp.Body.Bytes(), &body)
			Expect(err).ToNot(HaveOccurred())
			Expect(body).To(HaveLen(2))
			Expect(body[0].Name).To(Equal("Tetris"))
		})

		It("should return 500 if repository fails", func() {
			gr.On("ListGames").Return(nil, errors.New("db error"))
			resp := makeRequest(router, http.MethodGet, "/games", nil)
			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("POST /games", func() {
		It("should return 201 when game is created", func() {
			gr.On("GetGameByName", "Tetris").Return(nil, gorm.ErrRecordNotFound)
			gr.On("CreateGame", "Tetris").Return(&model.Game{ID: 1, Name: "Tetris"}, nil)

			payload := map[string]string{"name": "Tetris"}
			body, _ := json.Marshal(payload)

			resp := makeRequest(router, http.MethodPost, "/games", body)
			Expect(resp.Code).To(Equal(http.StatusCreated))

			var dto dto.GameDTO
			Expect(json.Unmarshal(resp.Body.Bytes(), &dto)).To(Succeed())
			Expect(dto.Name).To(Equal("Tetris"))
		})

		It("should return 409 when game already exists", func() {
			gr.On("GetGameByName", "Tetris").Return(&model.Game{ID: 1, Name: "Tetris"}, nil)

			body, _ := json.Marshal(map[string]string{"name": "Tetris"})
			resp := makeRequest(router, http.MethodPost, "/games", body)
			Expect(resp.Code).To(Equal(http.StatusConflict))
		})

		It("should return 500 if creation fails", func() {
			gr.On("GetGameByName", "Tetris").Return(nil, gorm.ErrRecordNotFound)
			gr.On("CreateGame", "Tetris").Return(nil, errors.New("insert error"))

			body, _ := json.Marshal(map[string]string{"name": "Tetris"})
			resp := makeRequest(router, http.MethodPost, "/games", body)
			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})

func makeRequest(r *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}

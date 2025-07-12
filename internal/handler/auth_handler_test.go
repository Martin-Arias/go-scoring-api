package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/Martin-Arias/go-scoring-api/internal/handler"
	"github.com/Martin-Arias/go-scoring-api/internal/mocks"
	"github.com/Martin-Arias/go-scoring-api/internal/model"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var _ = Describe("AuthHandler Register", func() {
	var (
		r        *gin.Engine
		mockRepo *mocks.UserRepositoryMock
		authH    *handler.AuthHandler
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		r = gin.Default()
		mockRepo = new(mocks.UserRepositoryMock)
		authH = handler.NewAuthHandler(mockRepo)
		r.POST("/register", authH.Register)
	})

	Context("when registration is successful", func() {
		It("should return 201 and create the user", func() {
			body := map[string]string{
				"username": "martin",
				"password": "secure123",
			}
			jsonBody, _ := json.Marshal(body)

			mockRepo.On("GetUserByUsername", "martin").Return(nil, nil)
			mockRepo.On("RegisterUser", mock.Anything).Return(nil)

			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusCreated))
			mockRepo.AssertExpectations(GinkgoT())
		})
	})

	Context("when registration fails", func() {
		It("should return 500 Internal server error", func() {
			body := map[string]string{
				"username": "martin",
				"password": "secure123",
			}
			jsonBody, _ := json.Marshal(body)

			mockRepo.On("GetUserByUsername", "martin").Return(nil, nil)
			mockRepo.On("RegisterUser", mock.Anything).Return(errors.New("error registering user"))

			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
			mockRepo.AssertExpectations(GinkgoT())
		})
	})

	Context("when username already exists", func() {
		It("should return 409 Conflict", func() {
			existing := &model.User{Username: "martin"}
			mockRepo.On("GetUserByUsername", "martin").Return(existing, nil)

			body := map[string]string{
				"username": "martin",
				"password": "secure123",
			}
			jsonBody, _ := json.Marshal(body)

			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusConflict))
		})
	})

	Context("when request is invalid", func() {
		It("should return 400 Bad request", func() {

			body := map[string]string{
				"username": "",
				"password": "",
			}
			jsonBody, _ := json.Marshal(body)

			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when user retrieval fails", func() {
		It("should return 500 Internal server error", func() {
			mockRepo.On("GetUserByUsername", "martin").Return(nil, gorm.ErrDuplicatedKey)

			body := map[string]string{
				"username": "martin",
				"password": "somepass",
			}
			jsonBody, _ := json.Marshal(body)

			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})

var _ = Describe("AuthHandler Login", func() {
	var (
		r        *gin.Engine
		mockRepo *mocks.UserRepositoryMock
		authH    *handler.AuthHandler
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		r = gin.Default()
		mockRepo = new(mocks.UserRepositoryMock)
		authH = handler.NewAuthHandler(mockRepo)
		r.POST("/login", authH.Login)
	})

	Context("when login is successful", func() {
		It("returns 200 and a token", func() {
			password := "secure123"
			hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			user := &model.User{
				ID:           1,
				Username:     "martin",
				PasswordHash: string(hash),
			}
			mockRepo.On("GetUserByUsername", "martin").Return(user, nil)

			payload := map[string]string{
				"username": "martin",
				"password": password,
			}
			body, _ := json.Marshal(payload)

			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Body.String()).To(ContainSubstring("token"))
			mockRepo.AssertExpectations(GinkgoT())
		})
	})

	Context("when user does not exist", func() {
		It("returns 401", func() {
			mockRepo.On("GetUserByUsername", "ghost").Return(nil, gorm.ErrRecordNotFound)

			body := map[string]string{
				"username": "ghost",
				"password": "anypass",
			}
			jsonBody, _ := json.Marshal(body)

			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("when password is incorrect", func() {
		It("returns 401", func() {
			hash, _ := bcrypt.GenerateFromPassword([]byte("realpass"), bcrypt.DefaultCost)
			user := &model.User{
				ID:           1,
				Username:     "martin",
				PasswordHash: string(hash),
			}
			mockRepo.On("GetUserByUsername", "martin").Return(user, nil)

			body := map[string]string{
				"username": "martin",
				"password": "wrongpass",
			}
			jsonBody, _ := json.Marshal(body)

			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("when request is invalid", func() {
		It("should return 400 Bad request", func() {

			body := map[string]string{
				"username": "",
				"password": "",
			}
			jsonBody, _ := json.Marshal(body)

			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when user retrieval fails", func() {
		It("should return 500 Internal server error", func() {
			mockRepo.On("GetUserByUsername", "martin").Return(nil, gorm.ErrDuplicatedKey)

			body := map[string]string{
				"username": "martin",
				"password": "somepass",
			}
			jsonBody, _ := json.Marshal(body)

			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})

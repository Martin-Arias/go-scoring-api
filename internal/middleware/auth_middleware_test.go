package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/Martin-Arias/go-scoring-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthMiddleware", func() {
	var r *gin.Engine
	var token string

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		r = gin.Default()

		r.GET("/protected", middleware.AuthMiddleware(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Authorized"})
		})

		// Genera un token válido
		secret := os.Getenv("JWT_SECRET")
		claims := jwt.MapClaims{
			"uid":      1,
			"username": "testuser",
			"admin":    false,
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(time.Hour).Unix(),
		}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		var err error
		token, err = t.SignedString([]byte(secret))
		Expect(err).To(BeNil())
	})

	Context("with valid token", func() {
		It("should return 200", func() {
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Body.String()).To(ContainSubstring("Authorized"))
		})
	})

	Context("with missing token", func() {
		It("should return 401", func() {
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusUnauthorized))
			Expect(resp.Body.String()).To(ContainSubstring("missing authorization header"))
		})
	})

	Context("with invalid token", func() {
		It("should return 401", func() {
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer invalidtoken123")
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusUnauthorized))
			Expect(resp.Body.String()).To(ContainSubstring("invalid access token"))
		})
	})

	Context("with wrong signing method", func() {
		It("should return 401", func() {
			// Generar un token con método incorrecto
			t := jwt.New(jwt.SigningMethodHS384) // No HMAC
			tokenStr, _ := t.SignedString([]byte("somekey"))

			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+tokenStr)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})

var _ = Describe("AdminMiddleware", func() {
	var r *gin.Engine

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		r = gin.Default()
	})

	Context("when user is admin", func() {
		It("allows access", func() {
			//setea el contexto como admin
			r.Use(func(c *gin.Context) {
				c.Set("admin", true)
				c.Next()
			})

			r.GET("/admin", middleware.AdminMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Welcome Admin!"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Body.String()).To(ContainSubstring("Welcome Admin!"))
		})
	})

	Context("when user is not admin", func() {
		It("returns forbidden", func() {
			//setea el contexto como no-admin
			r.Use(func(c *gin.Context) {
				c.Set("admin", false)
				c.Next()
			})

			r.GET("/admin", middleware.AdminMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Welcome Admin!"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusForbidden))
			Expect(resp.Body.String()).To(ContainSubstring("forbidden resource"))
		})
	})
})

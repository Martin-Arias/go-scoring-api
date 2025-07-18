package main

import (
	"log"
	"strconv"

	"github.com/Martin-Arias/go-scoring-api/cmd/api/handlers"
	_ "github.com/Martin-Arias/go-scoring-api/docs"
	"github.com/Martin-Arias/go-scoring-api/internal/middleware"
	repository "github.com/Martin-Arias/go-scoring-api/internal/repository/postgres"
	"github.com/Martin-Arias/go-scoring-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var db *gorm.DB

// Custom registry (without default Go metrics)
var customRegistry = prometheus.NewRegistry()

// Custom metrics handler with custom registry
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func setupRouter() *gin.Engine {

	sr := repository.NewScoreRepository(db)
	ur := repository.NewUserRepository(db)
	gr := repository.NewGameRepository(db)

	us := services.NewUserService(ur)
	ss := services.NewScoreService(sr, ur, gr)
	gs := services.NewGameService(gr)

	r := gin.Default()
	r.GET("/metrics", PrometheusHandler())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(RequestMetricsMiddleware())

	userHandler := handlers.NewUserHandler(us)
	gameHandler := handlers.NewGameHandler(gs)
	scoreHandler := handlers.NewScoreHandler(ss)
	// Public routes
	auth := r.Group("/auth")
	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	api.POST("/games", middleware.AdminMiddleware(), gameHandler.Create)
	api.GET("/games", gameHandler.List)

	api.PUT("/scores", middleware.AdminMiddleware(), scoreHandler.Submit)
	api.GET("/scores/user", scoreHandler.GetUserScores)
	api.GET("/scores/game", scoreHandler.GetGameScores)
	api.GET("/scores/game/stats", scoreHandler.GetGameStats)

	return r
}
func init() {
	customRegistry.MustRegister(HttpRequestTotal, HttpRequestErrorTotal)
	var err error
	db, err = repository.ConnectAndMigrate()
	if err != nil {
		log.Fatal("DB connection/migration failed:", err)
	}
}

// @title           Scoring API
// @version         1.0
// @description     API for managing players, games and scores
// @termsOfService  http://example.com/terms/

// @contact.name   Martín Arias
// @contact.email  martin@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

// Define metrics
var (
	HttpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_total",
		Help: "Total number of requests processed by the API",
	}, []string{"path", "status"})

	HttpRequestErrorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_error_total",
		Help: "Total number of errors returned by the API",
	}, []string{"path", "status"})
)

// Middleware to record incoming requests metrics
func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Next()
		status := c.Writer.Status()
		if status < 400 {
			HttpRequestTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
		} else {
			HttpRequestErrorTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
		}
	}
}

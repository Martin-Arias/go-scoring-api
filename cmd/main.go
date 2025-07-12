package main

import (
	"log"
	"os"

	"github.com/Martin-Arias/go-scoring-api/internal/handler"
	"github.com/Martin-Arias/go-scoring-api/internal/middleware"
	"github.com/Martin-Arias/go-scoring-api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	r := gin.Default()
	r.GET("/metrics", PrometheusHandler())
	r.Use(middleware.RequestMetricsMiddleware())

	authHandler := handler.NewAuthHandler(ur)
	gameHandler := handler.NewGameHandler(gr)
	scoreHandler := handler.NewScoreHandler(sr, ur, gr)
	// Public routes
	auth := r.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	api.POST("/games", middleware.AdminMiddleware(), gameHandler.Create)
	api.GET("/games", gameHandler.List)

	api.PUT("/scores", middleware.AdminMiddleware(), scoreHandler.Submit)
	api.GET("/scores/user", scoreHandler.GetScoresByPlayerID)
	api.GET("/scores/game", scoreHandler.GetScoresByGameID)
	api.GET("/scores/game/stats", scoreHandler.GetStatisticsByGameID)

	return r
}
func init() {
	customRegistry.MustRegister(middleware.HttpRequestTotal, middleware.HttpRequestErrorTotal)
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err = repository.ConnectAndMigrate(dsn)
	if err != nil {
		log.Fatal("DB connection/migration failed:", err)
	}
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

package main

import (
	"log"
	"strconv"

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
	r.Use(RequestMetricsMiddleware())

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
	customRegistry.MustRegister(HttpRequestTotal, HttpRequestErrorTotal)
	var err error
	db, err = repository.ConnectAndMigrate()
	if err != nil {
		log.Fatal("DB connection/migration failed:", err)
	}
}

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

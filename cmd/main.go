package main

import (
	"log"
	"os"

	"github.com/Martin-Arias/go-scoring-api/internal/handler"
	"github.com/Martin-Arias/go-scoring-api/internal/middleware"
	"github.com/Martin-Arias/go-scoring-api/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func setupRouter() *gin.Engine {

	sr := repository.NewScoreRepository(db)
	ur := repository.NewUserRepository(db)
	gr := repository.NewGameRepository(db)

	r := gin.Default()
	authHandler := handler.NewAuthHandler(db)
	gameHandler := handler.NewGameHandler(db)
	scoreHandler := handler.NewScoreHandler(sr, ur, gr)
	// Public routes
	auth := r.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	api.POST("/games", gameHandler.Create)
	api.GET("/games", gameHandler.List)
	api.PUT("/scores", scoreHandler.Submit)

	return r
}
func init() {
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

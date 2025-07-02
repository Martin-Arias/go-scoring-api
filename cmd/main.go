package main

import (
	"log"
	"os"

	"github.com/Martin-Arias/go-scoring-api/internal/handler"
	"github.com/Martin-Arias/go-scoring-api/internal/repository"
	"github.com/gin-gonic/gin" // Import the library
	"gorm.io/gorm"
)

var db *gorm.DB

func setupRouter() *gin.Engine {
	r := gin.Default()

	gameHandler := handler.NewGameHandler(db)
	//scoreHandler := handler.NewScoreHandler(db)

	// Admin routes
	r.POST("/games", gameHandler.Create)
	r.GET("/games", gameHandler.List)

	//r.POST("/games/score", scoreHandler.Submit)

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

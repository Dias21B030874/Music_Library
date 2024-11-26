package main

import (
	"log"
	"music_library/config"
	"music_library/database"
	_ "music_library/docs"
	"music_library/songs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Music Library API
// @version 1.0
// @description API для управления библиотекой песен.
// @host localhost:8080
// @BasePath /
func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("ERROR: Could not load environment variables: %v", err)
	}
	log.Println("INFO: Environment variables loaded.")

	db := database.Connect()
	if db == nil {
		log.Fatalf("ERROR: Could not connect to database")
	}
	log.Println("INFO: Database connection established.")
	database.Migrate(db)
	log.Println("INFO: Database migrations completed.")

	router := gin.Default()

	router.GET("/info", songs.GetSongVerses)
	router.GET("/songs", songs.GetSongs)
	router.GET("/songs/:id/verses", songs.GetSongVerses)
	router.PUT("/songs/:id", songs.UpdateSong)
	router.DELETE("/songs/:id", songs.DeleteSong)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("INFO: Swagger documentation is available at http://localhost:8080/swagger/index.html")

	log.Println("INFO: Starting the main server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("ERROR: Failed to start the server: %v", err)
	}
}

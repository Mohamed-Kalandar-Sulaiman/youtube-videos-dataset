package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/config"
	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/database"
	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/external"
	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/handlers"
	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/repository"
	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/services"

	fiber "github.com/gofiber/fiber/v2"
)

var (
	err error
)

func main() {
	// Loading env vars
	err = config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
	postgresConfig := &database.PostgresDB{
		Host:     config.Get("POSTGRES_HOST", "").(string),
		Port:     config.Get("POSTGRES_PORT", 5432).(int),
		User:     config.Get("POSTGRES_USER", "").(string),
		Password: config.Get("POSTGRES_PASSWORD", "").(string),
		Database: config.Get("POSTGRES_DB", "").(string),
	}

	db := database.GetDBInstance(postgresConfig)
	if db == nil {
		log.Fatal("Failed to establish a database connection")
	}
	defer database.CloseDB()

	videoRepo := repository.NewVideoRepository(db)

	// ! Run migrations
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	relativePath := "src/database/migrations/initialise.sql"
	sqlFilePath := filepath.Join(baseDir, relativePath)

	if err := videoRepo.RunMigrations(sqlFilePath); err != nil {
		log.Fatalf("Failed to execute migrations: %v", err)
	}

	videoService := services.NewVideoService(videoRepo)
	VideoHandler := handlers.NewVideoHandler(videoService)

	backGroundProcess := external.NewBackGroundProcess(config.Get("GOOGLE_API_KEY", "").(string), videoRepo)
	ctx, cancelBackGroundProcess := context.WithCancel(context.Background())
	go backGroundProcess.Start(ctx)
	defer cancelBackGroundProcess()

	//! Fiber app setup
	app := fiber.New()
	vidoes := app.Group("/api/v1")
	vidoes.Get("/videos", VideoHandler.GetVideos)

	if err != nil {
		log.Printf("Error checking table existence: %v", err)
	}

	log.Fatal(app.Listen(":3000"))

}

package main

import (
	"log"

	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/config"
	"github.com/Mohamed-Kalandar-Sulaiman/youtube-videos-dataset/src/database"

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
		Host:    config.Get("POSTGRES_HOST", "").(string),
		Port:     config.Get("POSTGRES_PORT", 5432).(int),
		User:     config.Get("POSTGRES_USER", "").(string),
		Password: config.Get("POSTGRES_PASSWORD", "").(string),
		Database: config.Get("POSTGRES_DB", "").(string),
	}

	db := database.GetDBInstance(postgresConfig)
	if db == nil {
		log.Fatal("Failed to establish a database connection")
	}

    log.Println(db.Stats())

	log.Println("Database connected successfully !!!")

	defer database.CloseDB()
    database.CloseDB()
	app := fiber.New()

	if err != nil {
		log.Printf("Error checking table existence: %v", err)
	}

	log.Fatal(app.Listen(":3000"))

}

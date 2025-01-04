package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return errors.New("failed to load env vars")
	}
	log.Println("Successfully loaded env vars")

	return nil
}

func Get(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	switch defaultValue.(type) {
	case int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Invalid int value for key %s: %v. Using default value %v", key, err, defaultValue)
			return defaultValue
		}
		return intValue
	case bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Invalid bool value for key %s: %v. Using default value %v", key, err, defaultValue)
			return defaultValue
		}
		return boolValue
	case string:
		return value
	default:
		log.Printf("Unsupported default value type for key %s. Returning default value.", key)
		return defaultValue
	}
}

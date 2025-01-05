package database

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

var (
	instance *sql.DB
	mutex    sync.Mutex
)

type Database interface {
	Connect() (*sql.DB, error)
	Close() error
}

func GetDBInstance(db Database) *sql.DB {
	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil {
		instance = createConnectionWithRetry(db, 3, time.Second*2)
	}
	return instance
}

func createConnectionWithRetry(db Database, retries int, backoff time.Duration) *sql.DB {
	var conn *sql.DB
	var err error

	for i := 0; i <= retries; i++ {
		log.Printf("Attempt %d to connect to the database\n", i+1)

		conn, err = db.Connect()
		if err == nil {
			log.Println("Database connection established successfully")
			return conn
		}

		log.Printf("Database connection failed: %v", err)
		if i < retries {
			log.Printf("Retrying in %v...\n", backoff)
			time.Sleep(backoff)
			backoff *= 2
		} else {
			log.Fatalf("Failed to establish database connection after %d attempts", retries)
		}
	}
	return nil
}

func CloseDB() {
	mutex.Lock()
	defer mutex.Unlock()

	if instance != nil {
		log.Println("Closing database connection")
		if err := instance.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully")
			instance = nil
		}
	}
}

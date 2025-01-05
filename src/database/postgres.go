package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (p *PostgresDB) Connect() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.Database,
	)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL database successfully")
	return db, nil
}

func (p *PostgresDB) Close() error {
	if instance != nil {
		log.Println("Closing PostgreSQL database connection")
		err := instance.Close()
		if err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
		log.Println("PostgreSQL database connection closed successfully")
		instance = nil
		return nil
	}
	log.Println("No PostgreSQL database connection to close")
	return nil
}



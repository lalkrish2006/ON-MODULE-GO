package database

import (
	"database/sql"
	"log"
	"od-system/internal/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Connect establishes the database connection
func Connect(cfg config.Config) {
	var err error
	DB, err = sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to database successfully")
}

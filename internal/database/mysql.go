package database

import (
	"database/sql"
	"fmt"
	"healing_photons/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

// InitializeDB establishes a connection to the MySQL database
func InitializeDB(cfg *config.Config) (*sql.DB, error) {
	// Construct connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&tls=%s&parseTime=true",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.Port, cfg.DBName, cfg.UseSSL)

	// Open database connection
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Verify database connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var sqlDB *sql.DB

func ConnectDB(ctx context.Context) error {

	LoadEnv()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	newLogger := logger.New(
		log.New(os.Stdout, "[GORM] ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             300 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 250,
		Logger:          newLogger,
	})

	if err != nil {
		return fmt.Errorf("Failed to connect to database: %w", err)
	}

	sqlDB, err = database.DB()
	if err != nil {
		return fmt.Errorf("get sql db failed: %w", err)
	}

	// Enforce connection timeout
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("db ping timeout: %w", err)
	}

	DB = database
	log.Println("Database connected")
	return nil
}

func CloseSqlConnection() {
	if err := sqlDB.Close(); err != nil {
		log.Println("Error in closing sql connection", err)
	}
}

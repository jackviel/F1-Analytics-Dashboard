package config

import (
	"fmt"
	"os"

	"github.com/f1-analytics/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	// Get database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create DSN string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set to logger.Silent in production
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(3600) // 1 hour

	// Store the database connection
	DB = db

	return nil
}

// AutoMigrate performs database migrations
func AutoMigrate() error {
	// Import all models here
	models := []interface{}{
		&models.Driver{},
		&models.Team{},
		&models.Race{},
		&models.Circuit{},
		&models.RaceDriver{},
		&models.RaceTeam{},
		&models.Lap{},
	}

	// Run migrations
	if err := DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}

// GetDB returns the global database instance
func GetDB() *gorm.DB {
	return DB
}

package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupDatabase initializes the database connection with GORM
func SetupDatabase(databaseURL string, env string) (*gorm.DB, error) {
	// Configure GORM logger based on environment
	logLevel := logger.Info
	if env == "production" {
		logLevel = logger.Warn
	}

	config := &gorm.Config{
		Logger:      logger.Default.LogMode(logLevel),
		PrepareStmt: false, // Disable prepared statements for Supabase connection pooling
	}

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(databaseURL), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return db, nil
}

// AutoMigrate runs automatic migrations for all models
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}
	log.Println("✅ Database migrations completed")
	return nil
}

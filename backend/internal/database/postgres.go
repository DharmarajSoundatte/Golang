package database

import (
	"fmt"
	"time"

	"github.com/DharmarajSoundatte/Golang/backend/internal/config"
	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgres opens a new PostgreSQL connection, verifies it with a ping,
// configures the pool, and runs auto-migrations.
func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// ── Verify connectivity ───────────────────────────────────────────────────
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed — check credentials/network/CA cert: %w", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&models.User{},
	); err != nil {
		return nil, fmt.Errorf("auto-migrate: %w", err)
	}

	return db, nil
}

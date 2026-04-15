package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DharmarajSoundatte/Golang/backend/internal/config"
	"github.com/DharmarajSoundatte/Golang/backend/internal/database"
	"github.com/DharmarajSoundatte/Golang/backend/internal/routes"
	"go.uber.org/zap"
)

func main() {
	// ── Logger ────────────────────────────────────────────────────────────────
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := database.NewPostgres(cfg)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	logger.Info("database connected successfully")

	// ── Router ────────────────────────────────────────────────────────────────
	r := routes.Setup(db, cfg, logger)

	// ── HTTP Server ───────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── Graceful Shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("forced shutdown", zap.Error(err))
	}
	logger.Info("server exited gracefully")
}

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"mewsoproxy/server/config"
	"mewsoproxy/server/database"
	"mewsoproxy/server/pkg/redis"
	"mewsoproxy/server/router"
)

func main() {
	cfg := config.Load()
	if err := config.EnsureSecrets(cfg, cfg.App.SecretsFile); err != nil {
		slog.Error("ensure secrets failed", "err", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		slog.Error("config invalid", "err", err)
		os.Exit(1)
	}

	db, err := database.InitDB(cfg.Database)
	if err != nil {
		slog.Error("database init failed", "err", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("database pool failed", "err", err)
		os.Exit(1)
	}
	sqlDB.SetMaxOpenConns(16)
	sqlDB.SetMaxIdleConns(8)
	sqlDB.SetConnMaxLifetime(time.Hour)
	defer sqlDB.Close()

	bctx := context.Background()
	if err := database.ApplySystemConfig(bctx, db, cfg); err != nil {
		slog.Warn("apply system config failed", "err", err)
	}
	if err := database.SeedAdmin(bctx, db, cfg); err != nil {
		slog.Error("seed admin failed", "err", err)
		os.Exit(1)
	}

	rds := redis.New(cfg.Redis)
	ctx := context.Background()
	if err := rds.Ping(ctx); err != nil {
		slog.Error("redis ping failed", "err", err)
		os.Exit(1)
	}
	defer rds.Close()

	engine := router.New(cfg, db, rds)
	srv := &http.Server{
		Addr:    addr(cfg.Server.Port),
		Handler: engine,
	}

	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server listen failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "err", err)
	}
}

func addr(port int) string {
	return ":" + strconv.Itoa(port)
}

package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/magneless/todo-app/internal/config"
	"github.com/magneless/todo-app/internal/http-server/router"
	"github.com/magneless/todo-app/internal/lib/logger/sl"
	repository "github.com/magneless/todo-app/internal/repository"
	"github.com/magneless/todo-app/internal/storage/postgre"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting todo-app", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgre.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	repo := repository.New(storage)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router.New(log, repo),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {

	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

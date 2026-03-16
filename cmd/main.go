package main

import (
	"context"
	"habits/internal/config"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// TODO:
	// initialize db, [done]
	// migrations, [done]
	// habit repo,
	// users repo,
	// report repo,
	// habit usecase,
	// users usecase,
	// report usecase,
	// create api handler(private/public),
	// start http server
	cfg, err := config.Load()
	if err != nil {
		slog.Error("error", slog.Any("err", err))
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		slog.Error("error connecting to db", "error", err)
		return
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		slog.Error("error pinging db", "error", err)
		return
	}

	slog.Info("connected to db successfully")

	m, err := migrate.New(
		"file://migrations",
		cfg.DB.DSN,
	)
	if err != nil {
		slog.Error("error creating migrate instance", "error", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("error applying migrations", "error", err)
		return
	}

	slog.Info("migrations applied successfully")

	
}

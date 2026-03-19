package main

import (
	"context"
	"habits/internal/config"
	"habits/internal/habits"
	http "habits/internal/http/public"
	"habits/internal/postgres"
	"habits/internal/users"
	"log/slog"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

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

	r := gin.Default()

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		slog.Error("error creating oidc provider", "error", err)
		return
	}
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,	
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL: cfg.Google.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.Google.ClientID,
	})

	slog.Info("migrations applied successfully")

	habitsRepo := postgres.NewHabitsRepository(pool)
	usersRepo := postgres.NewUserRepository(pool)
	habitUsecase := habits.NewHabitUsecase(habitsRepo)
	userUsecase := users.NewUsersUsecase(usersRepo)
	handlers := http.NewUserHandlers(userUsecase, habitUsecase)

	r.Handle("GET", "/health", func (c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.Handle("POST", "/api/auth/google/login", func(ctx *gin.Context) {
		handlers.LoginWithGoogle(ctx, oauthConfig)
	})
	r.Handle("POST", "/api/auth/google/callback", func(ctx *gin.Context) {
		handlers.CallbackWithGoogle(ctx, oauthConfig, verifier)
	})
	// r.Handle("POST", "/api/auth/google/logout", handlers.LoginWithGoogle)
	// r.Handle("POST", "/api/auth/me", handlers.GetUser)
	r.Run(":8080")
	
}

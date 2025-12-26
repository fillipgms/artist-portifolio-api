package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/fillipgms/portfolio-api/internal/env"
	"github.com/fillipgms/portfolio-api/internal/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main () {
	godotenv.Load()

	ctx := context.Background()

	cfg := config{
		address: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=portfolio sslmode=disable"),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("Connected to Database", "dsn", cfg.db.dsn)

	helpers.ConnectToBunny()

	api := &application{
		config: cfg,
		db: conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server Failed to Start", "error", err)
		os.Exit(1)
	}
}
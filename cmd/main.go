package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/ybuilds/ecomm-api/internal/env"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("error loading .env file", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	cfg := config{
		addr: ":8000",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=kjm40438 dbname=ecom sslmode=disable"),
		},
	}

	api := api{
		config: cfg,
	}

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error("connection to database failed", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	logger.Info("connected to database", "dsn", cfg.db.dsn)

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}
}

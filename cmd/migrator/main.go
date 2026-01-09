package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	var (
		migrationsDir = os.Getenv("DIR")
		dsn           = os.Getenv("DSN")
		driverName    = "postgres"
	)

	flag.Parse()

	log := setupLogger()

	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		log.Error("error opening database connection", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err = db.Ping(); err != nil {
		log.Error("error pinging database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Error("command is required: up, down, status")
		os.Exit(1)
	}

	command := args[0]

	if err = goose.RunContext(context.Background(), command, db.DB, migrationsDir, args[1:]...); err != nil {
		log.Error("error running migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("migrations complete", slog.String("command", command))

}

// setupLogger инициализирует логгер с JSON-обработчиком
func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}

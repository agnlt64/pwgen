package main

import (
	"context"
	"log"
	"os"
	"pwgen/internal/commands"
	"pwgen/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func connect(ctx context.Context, dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dbURL)

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	return pool
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Args) < 2 {
		log.Fatal("not enough args")
	}

	dbURL    := os.Getenv("DB_URL")
	pool     := connect(context.Background(), dbURL)
	queries  := queries.NewQueries(pool)
	commands := commands.NewCommands(queries)

	subCmd := os.Args[1]
	switch subCmd {
	case "new-vault":
		commands.NewVault()
	case "new-pass":
		commands.NewPass()
	case "get-pass":
		commands.GetPass()
	// todo: help subcommand
	default:
		log.Fatalf("%s is not a valid subcommand", subCmd)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	c "pwgen/internal/commands"
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
		fmt.Println("Not enough arguments. See usage below.")
		c.Usage()
		os.Exit(1)
	}

	dbURL    := os.Getenv("DB_URL")
	pool     := connect(context.Background(), dbURL)
	queries  := queries.NewQueries(pool)
	commands := c.NewCommands(queries, os.Args[1:])

	subCmd := os.Args[1]
	switch subCmd {
	case "new-vault":
		commands.NewVault()
	case "new-pass":
		commands.NewPass()
	case "get-pass":
		commands.GetPass()
	case "help":
		c.Usage()
	default:
		fmt.Printf("%s is not a valid command. See usage below.\n", subCmd)
		c.Usage()
		os.Exit(1)
	}
}

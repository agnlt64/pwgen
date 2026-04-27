package cmd

import (
	"context"
	"os"
	"log"

	"pwgen/internal/db"
	"github.com/spf13/cobra"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var queries *db.Queries

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pwgen",
	Short: "A CLI to securely manage your passwords",
	Args: cobra.MinimumNArgs(1),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func connect(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	return pool, err
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		pool, err := connect(context.Background(), os.Getenv("DB_URL"))
		queries = db.New(pool)
		return err
	}
}

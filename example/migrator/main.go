package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/pkg/migrator"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}
	},
}

func init() {
	db := postgres.New(postgres.Config{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "password",
		DBName:   "postgres",
		SSLMode:  "disable",
	})

	migrator := migrator.NewMigrator(db.Conn(), embedMigrations, migrator.Config{
		Dialect: "postgres",
		Dir:     "migrations",
	})

	cmd := migrator.UpCommand()

	rootCmd.AddCommand(cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

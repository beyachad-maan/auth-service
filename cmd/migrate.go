package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beyachad-maan/auth-service/pkg/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const postgresHost = "postgres-auth-service"

var migrateCmd = cobra.Command{
	Use:  "migrate",
	Run:  migrateDatabase,
	Long: "Run database migrations to ensure database is up to date.",
}

var config postgres.Config

func init() {
	// Add flags to the migrate command.
	migrateCmdFlags := migrateCmd.Flags()
	migrateCmdFlags.String("database-name", "anglit-db", "Database name")
	migrateCmdFlags.Int("database-port", 5432, "Database port")
	migrateCmdFlags.String("database-user", "postgres", "Database user")
	migrateCmdFlags.String("database-password", "postgres", "Database password")

	// Bind the flags to the viper configuration.
	viper.BindPFlag("database-name", migrateCmd.Flags().Lookup("database-name"))
	viper.BindPFlag("database-port", migrateCmd.Flags().Lookup("database-port"))
	viper.BindPFlag("database-user", migrateCmd.Flags().Lookup("database-user"))
	viper.BindPFlag("database-password", migrateCmd.Flags().Lookup("database-password"))

	// Set default values for the database configuration environment variables.
	// This is crucial for the viper to consider the environment variables.
	viper.BindEnv("DATABASE_USER")
	viper.SetDefault("DATABASE_USER", "postgres")
	viper.BindEnv("DATABASE_PASSWORD")
	viper.SetDefault("DATABASE_PASSWORD", "postgres")
	viper.BindEnv("DATABASE_NAME")
	viper.SetDefault("DATABASE_NAME", "anglit-db")
	viper.BindEnv("DATABASE_PORT")
	viper.SetDefault("DATABASE_PORT", 5432)

	viper.AutomaticEnv()
	viper.Unmarshal(&config)
}

func migrateDatabase(cmd *cobra.Command, args []string) {
	log.Default().Printf("Migrating database...")
	log.Default().Printf("Recieved config: %+v", config)
	for attempt := 0; attempt < 3; attempt++ {
		m, err := migrate.New(
			"file://db/migrations",
			fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
				config.DatabaseUser,
				config.DatabasePassword,
				postgresHost,
				config.DatabasePort,
				config.DatabaseName,
			),
		)

		if err != nil {
			log.Printf("Error occurred attempting to connect to db instance: %v. Retrying...", err)
			time.Sleep(time.Second * 10)
			continue
		}
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
		log.Default().Printf("Database migration finished successfully.")
		return
	}

	log.Default().Printf("Database migration Failed.")
	os.Exit(1)
}

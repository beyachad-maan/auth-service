package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const postgresDriver = "postgres"

type Config struct {
	DatabaseName     string `mapstructure:"DATABASE_NAME"`
	DatabasePort     int    `mapstructure:"DATABASE_PORT"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
}

func ConnectDB(user string, password string, host string, port int, dbName string) (*sqlx.DB, error) {
	return sqlx.Connect(postgresDriver, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbName,
	))
}

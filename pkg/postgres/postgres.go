package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type Config struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
	Database string `yaml:"database" env:"POSTGRES_DB" env-required:"true"`
}

func New(ctx context.Context, config *Config) (*pgx.Conn, error) {
	slog.Info("Connecting to PostgresSQL...")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		slog.Error("Failed to connect to PostgresSQL", err)
	}
	return conn, nil
}

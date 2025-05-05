package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"
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
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	err = Migrate(config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Migrate(cfg *Config) error {
	goose.SetVerbose(true)
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)
	sqlDB, err := goose.OpenDBWithDriver("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("goose open: %w", err)
	}
	defer sqlDB.Close()

	if err := goose.Up(sqlDB, "./db/migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

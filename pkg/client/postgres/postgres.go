package postgres

import (
	"context"
	"evo-test/internal/config"
	"github.com/jackc/pgx"
)

func NewClient(cfg *config.Repository) (*pgx.Conn, error) {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.Username,
		Password: cfg.Password,
	})
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

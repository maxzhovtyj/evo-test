package postgres

import (
	"context"
	"evo-test/internal/config"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func NewClient(cfg *config.Repository) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

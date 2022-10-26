package repository

import "github.com/jackc/pgx"

type Repository interface {
}

type repository struct {
	db *pgx.Conn
}

func New(conn *pgx.Conn) Repository {
	return &repository{db: conn}
}

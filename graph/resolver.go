package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/jackc/pgx/v4"
)

type Resolver struct {
	db *pgx.Conn
}

func NewResolver(db *pgx.Conn) *Resolver {
	return &Resolver{db: db}
}

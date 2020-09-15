package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/sql-judge/api-graphql/config"
)

func ConnectWithConfig(cfg config.PostgresConfig) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), cfg.ConnectionString())
}

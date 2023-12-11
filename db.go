package pgxmigrations

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row

	Begin(context.Context) (pgx.Tx, error)
}

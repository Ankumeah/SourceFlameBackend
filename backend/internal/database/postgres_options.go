//go:build postgres

package database

import (
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"errors"
)

const driver_name = "pgx"

func is_unique_violation(err error) bool {
	var pg_err *pgconn.PgError
	return errors.As(err, &pg_err) && pg_err.Code == "23505"
}

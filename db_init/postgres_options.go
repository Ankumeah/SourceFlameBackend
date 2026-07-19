//go:build postgres

package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun/dialect/pgdialect"
)

const driverName = "pgx"

var dia = pgdialect.New()

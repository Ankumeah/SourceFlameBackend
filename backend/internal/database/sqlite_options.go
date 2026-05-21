//go:build sqlite3
package database

import (
  _ "github.com/mattn/go-sqlite3"
  "github.com/uptrace/bun/dialect/sqlitedialect"
)

const driver_name = "sqlite3"
var dia = sqlitedialect.New()

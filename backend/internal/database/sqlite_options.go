//go:build sqlite3
package database

import "github.com/mattn/go-sqlite3"

const driver_name = "sqlite3"
func is_unique_violation(err error) bool {
	sqlite_err, ok := err.(sqlite3.Error)
  if !ok { return false }
	return sqlite_err.ExtendedCode == sqlite3.ErrConstraintUnique
}

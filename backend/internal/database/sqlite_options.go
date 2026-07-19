//go:build sqlite3

package database

import "github.com/mattn/go-sqlite3"

const driverName = "sqlite3"

func isUniqueViolation(err error) bool {
	sqliteErr, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}
	return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}

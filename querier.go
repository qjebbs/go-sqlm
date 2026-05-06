package sqlm

import (
	"database/sql"
)

// Querier is the interface for *sql.DB, *sql.Tx, etc.
type Querier interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

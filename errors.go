package sqlm

import "errors"

// errors
var (
	ErrNilDB        = errors.New("db is nil")
	ErrMultipleRows = errors.New("sql: multiple rows in result set")
)

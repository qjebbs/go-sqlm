package sqlm

import (
	"context"
	"database/sql"

	"github.com/qjebbs/go-sqlb"
	"github.com/qjebbs/go-sqlf/v4"
	"github.com/qjebbs/go-sqlm/option"
)

// QueryOne is like Query but expects exactly one result.
// It returns sql.ErrNoRows if there are no results, and an error if there are more than one result.
func QueryOne[T any](ctx sqlb.Context, db Querier, b sqlf.Builder, fn func() (T, []any), options ...option.Option) (T, error) {
	r, err := scanBuilder(ctx, "QueryOne", db, b, fn, options...)
	if err != nil {
		var zero T
		return zero, err
	}
	if len(r) == 0 {
		var zero T
		return zero, sql.ErrNoRows
	}
	if len(r) > 1 {
		var zero T
		return zero, ErrMultipleRows
	}
	return r[0], nil
}

// Query executes a query from builder and scans the results using a provider function.
// The provider fn is called for each row to get the destination value and scan fields.
//
// Unlike Select, it doesn't require the builder to be a SelectBuilder.
// It's useful when you want to execute a raw query or a non-select query (like INSERT..RETURNING) with scanning results.
func Query[T any](ctx sqlb.Context, db Querier, b sqlf.Builder, fn func() (T, []any), options ...option.Option) ([]T, error) {
	return scanBuilder(ctx, "Query", db, b, fn, options...)
}

func scanBuilder[T any](ctx sqlb.Context, name string, db Querier, b sqlf.Builder, fn func() (T, []any), options ...option.Option) ([]T, error) {
	opt := option.New(options...)
	var debugger *debugger
	if opt.Debug.Enabled {
		value, _ := fn()
		debugger = newDebugger(name, value, opt)
		defer debugger.print(ctx.BaseDialect())
	}
	query, args, err := sqlf.Build(ctx, b)
	if err != nil {
		return nil, err
	}
	if debugger != nil {
		debugger.onBuilt(query, args)
	}
	if db == nil {
		return nil, ErrNilDB
	}
	r, err := scanQuery(ctx, db, query, args, debugger, fn)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// scanQuery scans query rows with scanner
func scanQuery[T any](ctx context.Context, db Querier, query string, args []any, debugger *debugger, fn func() (T, []any)) ([]T, error) {
	rows, err := db.Query(query, args...)
	if debugger != nil {
		debugger.onExec(err)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		dest, fields := fn()
		err = rows.Scan(fields...)
		if err != nil {
			return nil, err
		}
		results = append(results, dest)
	}
	if debugger != nil {
		debugger.onScan(len(results), err)
	}
	return results, rows.Err()
}

// // ScanRow scans a single row to dest, unlike rows.Scan(), it drops the extra columns.
// // It's useful when *sqlb.SelectBuilder.OrderBy() add extra column to the query.
// func scanRow(rows *sql.Rows, dest ...any) error {
// 	cols, err := rows.Columns()
// 	if err != nil {
// 		return err
// 	}
// 	nBlacholes := len(cols) - len(dest)
// 	for i := 0; i < nBlacholes; i++ {
// 		dest = append(dest, Blackhole)
// 	}
// 	return rows.Scan(dest...)
// }

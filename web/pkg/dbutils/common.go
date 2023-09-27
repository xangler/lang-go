package dbutils

import (
	"context"
	"database/sql"
)

type LockType int

const (
	NoLock LockType = iota
	WriteLock
)

type Scanner interface {
	Scan(...interface{}) error
}

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

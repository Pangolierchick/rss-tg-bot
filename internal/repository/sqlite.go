package repository

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"
	"time"
)

//go:embed schema.sql
var schemaFS embed.FS

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Init(ctx context.Context) error {
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, string(schema)); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	return nil
}

func (r *Repository) Begin(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func unixTime(t time.Time) int64 {
	return t.UTC().Unix()
}

func unixTimePtr(t *time.Time) any {
	if t == nil {
		return nil
	}

	return unixTime(*t)
}

func scanUnixTime(v int64) time.Time {
	return time.Unix(v, 0).UTC()
}

func scanNullUnixTime(v sql.NullInt64) *time.Time {
	if !v.Valid {
		return nil
	}

	t := scanUnixTime(v.Int64)
	return &t
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}

	return strings.TrimRight(strings.Repeat("?,", n), ",")
}

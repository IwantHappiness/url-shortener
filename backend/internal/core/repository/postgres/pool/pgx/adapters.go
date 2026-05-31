package core_pgx_pool

import (
	"errors"
	"fmt"

	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxRows struct {
	pgx.Rows
}

type PgxRow struct {
	pgx.Row
}

func (r *PgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return MapErrors(err)
	}

	return nil
}

type PgxCommandTag struct {
	pgconn.CommandTag
}

func MapErrors(err error) error {
	const pgxViolatesForeignKeyErrorCode = "23503"
	const pgxUniqueViolationErrorCode = "23505"

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code == pgxViolatesForeignKeyErrorCode {
			return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrViolatesForeignKey)
		}

		if pgErr.Code == pgxUniqueViolationErrorCode {
			return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrUniqueViolation)
		}
	}

	return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrUnknown)
}

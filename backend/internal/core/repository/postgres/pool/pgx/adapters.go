package core_pgx_pool

import (
	"errors"

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
		if errors.Is(err, pgx.ErrNoRows) {
			return core_postgres_pool.ErrNoRows
		}

		return err
	}

	return nil
}

type PgxCommandTag struct {
	pgconn.CommandTag
}

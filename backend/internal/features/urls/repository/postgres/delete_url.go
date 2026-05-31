package url_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
)

func (r *URLRepository) DeleteURL(ctx context.Context, urlID int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM url_shortener.urls WHERE id = $1`

	cmdTag, err := r.pool.Exec(ctx, query, urlID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("url with id='%d': %w", urlID, core_errors.ErrNotFound)
	}

	return nil
}

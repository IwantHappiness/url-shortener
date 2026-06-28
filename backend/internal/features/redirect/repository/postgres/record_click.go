package redirect_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
)

func (r *RedirectRepository) RecordClick(ctx context.Context, shortURL, ip string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO url_shortener.clicks
	(short_url, ip) VALUES ($1, $2);
	`

	_, err := r.pool.Exec(ctx, query, shortURL, ip)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return fmt.Errorf("url with short_url='%s': %w", shortURL, core_errors.ErrNotFound)
		}
		return fmt.Errorf("insert click: %w", err)
	}

	return nil
}

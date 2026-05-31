package redirect_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
)

func (r *RedirectRepository) GetByShortURL(ctx context.Context, shortURL string) (domain.Link, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	query := `SELECT original_url, short_url FROM url_shortener.urls WHERE short_url = $1`

	row := r.pool.QueryRow(ctx, query, shortURL)

	var model RedirectModel
	err := row.Scan(&model.OriginalURL, &model.ShortURL)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Link{}, fmt.Errorf("url with short_url='%s': %w", shortURL, core_errors.ErrNotFound)
		}
		return domain.Link{}, fmt.Errorf("scan error: %w", err)
	}

	linkDomain := redirectDomainFromModel(model)

	return linkDomain, nil
}

package url_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
)

func (r *URLRepository) CreateURL(ctx context.Context, link domain.Link) (domain.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO url_shortener.urls
	(user_id, original_url, short_url) VALUES ($1, $2, $3)
	RETURNING id, version, user_id, original_url, short_url, created_at;
	`

	row := r.pool.QueryRow(ctx, query, link.UserID, link.OriginalURL, link.ShortURL)

	var linkModel LinkModel
	err := row.Scan(
		&linkModel.ID,
		&linkModel.Version,
		&linkModel.UserID,
		&linkModel.OriginalURL,
		&linkModel.ShortURL,
		&linkModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Link{}, fmt.Errorf("%v: user with id='%d': %w", err, link.UserID, core_errors.ErrNotFound)
		}

		if errors.Is(err, core_postgres_pool.ErrUniqueViolation) {
			return domain.Link{}, fmt.Errorf("%v: duplicate value: %w", err, core_errors.ErrConflict)
		}

		return domain.Link{}, fmt.Errorf("scan error: %w", err)
	}

	linkDomain := urlDomainFromModel(linkModel)
	return linkDomain, nil
}

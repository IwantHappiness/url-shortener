package url_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
)

func (r *URLRepository) PatchURL(ctx context.Context, urlID int, link domain.Link) (domain.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE url_shortener.urls
	SET short_url = $1, version = version + 1
	WHERE id = $2 AND version = $3
	RETURNING id, version, user_id, short_url, original_url, created_at`

	row := r.pool.QueryRow(ctx, query, link.ShortURL, urlID, link.Version)

	var urlModel LinkModel
	err := row.Scan(&urlModel.ID, &urlModel.Version, &urlModel.UserID, &urlModel.ShortURL, &urlModel.OriginalURL, &urlModel.CreatedAt)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Link{}, fmt.Errorf("url with id='%d' concurrently accessed: %w", urlID, core_errors.ErrConflict)
		}

		return domain.Link{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := urlDomainFromModel(urlModel)
	return userDomain, nil
}

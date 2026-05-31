package url_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
)

func (r *URLRepository) GetURL(ctx context.Context, id int) (domain.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT id, version, user_id, original_url, short_url, created_at FROM url_shortener.urls WHERE id = $1;`

	row := r.pool.QueryRow(ctx, query, id)

	var linkModel LinkModel

	err := row.Scan(&linkModel.ID, &linkModel.Version, &linkModel.UserID, &linkModel.OriginalURL, &linkModel.ShortURL, &linkModel.CreatedAt)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Link{}, fmt.Errorf("url with id='%d' not found: %w", id, core_errors.ErrNotFound)
		}
		return domain.Link{}, fmt.Errorf("scan error: %w", err)
	}

	linkDomain := urlDomainFromModel(linkModel)
	return linkDomain, nil

}

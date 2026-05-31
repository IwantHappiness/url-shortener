package url_postgres_repository

import (
	"context"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

func (r *URLRepository) GetURLs(ctx context.Context, userID, limit, offset *int) ([]domain.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, user_id, original_url, short_url, created_at
	FROM url_shortener.urls
	%s
	ORDER BY id ASC
	LIMIT $1 OFFSET $2`

	args := []any{limit, offset}

	if userID != nil {
		query = fmt.Sprintf(query, "WHERE user_id=$3")
		args = append(args, userID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select urls: %w", err)
	}
	defer rows.Close()

	var urlsModel []LinkModel
	for rows.Next() {
		var link LinkModel
		if err := rows.Scan(&link.ID, &link.Version, &link.UserID, &link.OriginalURL, &link.ShortURL, &link.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan url: %w", err)
		}
		urlsModel = append(urlsModel, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	urlsDomain := ulrsDomainsFromModels(urlsModel)

	return urlsDomain, nil
}

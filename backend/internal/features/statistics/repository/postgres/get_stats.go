package stats_postgres_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
)

func (r *StatsRepository) GetStats(ctx context.Context, shortURL string, from, to *time.Time) (domain.LinkStats, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT
		u.short_url,
		u.original_url,
		u.created_at,
		COUNT(c.id) AS total_clicks,
		COUNT(DISTINCT c.ip) AS unique_ips,
		MAX(c.clicked_at) AS last_clicked_at
	FROM url_shortener.urls u
	LEFT JOIN url_shortener.clicks c
		ON c.short_url = u.short_url
		AND ($2::timestamptz IS NULL OR c.clicked_at >= $2)
		AND ($3::timestamptz IS NULL OR c.clicked_at <= $3)
	WHERE u.short_url = $1
	GROUP BY u.short_url, u.original_url, u.created_at;
	`

	row := r.pool.QueryRow(ctx, query, shortURL, from, to)

	var (
		model         StatsModel
		lastClickedAt sql.NullTime
	)

	err := row.Scan(
		&model.ShortURL,
		&model.OriginalURL,
		&model.CreatedAt,
		&model.TotalClicks,
		&model.UniqueIPs,
		&lastClickedAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.LinkStats{}, fmt.Errorf("url with short_url='%s': %w", shortURL, core_errors.ErrNotFound)
		}
		return domain.LinkStats{}, fmt.Errorf("select stats: %w", err)
	}

	if lastClickedAt.Valid {
		clickedAt := lastClickedAt.Time
		model.LastClickedAt = &clickedAt
	}

	return statsDomainFromModel(model), nil
}

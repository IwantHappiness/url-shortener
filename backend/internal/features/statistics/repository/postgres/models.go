package stats_postgres_repository

import (
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type StatsModel struct {
	ShortURL      string
	OriginalURL   string
	CreatedAt     time.Time
	TotalClicks   int64
	UniqueIPs     int64
	LastClickedAt *time.Time
}

func statsDomainFromModel(model StatsModel) domain.LinkStats {
	return domain.LinkStats{
		ShortURL:      model.ShortURL,
		OriginalURL:   model.OriginalURL,
		CreatedAt:     model.CreatedAt,
		TotalClicks:   int(model.TotalClicks),
		UniqueIPs:     int(model.UniqueIPs),
		LastClickedAt: model.LastClickedAt,
	}
}

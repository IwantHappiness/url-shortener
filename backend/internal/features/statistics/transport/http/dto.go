package stats_transport_http

import (
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type StatsDTOResponse struct {
	ShortURL      string     `json:"short_url"`
	OriginalURL   string     `json:"original_url"`
	CreatedAt     time.Time  `json:"created_at"`
	TotalClicks   int        `json:"total_clicks"`
	UniqueIPs     int        `json:"unique_ips"`
	LastClickedAt *time.Time `json:"last_clicked_at,omitempty"`
}

func statsDTOFromDomain(stats domain.LinkStats) StatsDTOResponse {
	return StatsDTOResponse{
		ShortURL:      stats.ShortURL,
		OriginalURL:   stats.OriginalURL,
		CreatedAt:     stats.CreatedAt,
		TotalClicks:   stats.TotalClicks,
		UniqueIPs:     stats.UniqueIPs,
		LastClickedAt: stats.LastClickedAt,
	}
}

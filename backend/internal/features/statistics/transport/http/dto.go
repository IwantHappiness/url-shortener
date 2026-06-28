package stats_transport_http

import (
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type StatsDTOResponse struct {
	ShortURL      string     `json:"short_url"             example:"abc123"`
	OriginalURL   string     `json:"original_url"          example:"https://example.com"`
	CreatedAt     time.Time  `json:"created_at"            example:"2024-06-28T15:04:05Z07:00"`
	TotalClicks   int        `json:"total_clicks"          example:"42"`
	UniqueIPs     int        `json:"unique_ips"            example:"15"`
	LastClickedAt *time.Time `json:"last_clicked_at,omitempty" example:"2024-06-28T15:04:05Z07:00"`
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

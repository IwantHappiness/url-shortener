package stats_transport_http

import (
	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	http_types "github.com/IwantHappiness/url-shortener/internal/core/transport/http/types"
)

type StatsDTOResponse struct {
	ShortURL      string               `json:"short_url"             example:"abc123"`
	OriginalURL   string               `json:"original_url"          example:"https://example.com"`
	CreatedAt     http_types.DateOnly  `json:"created_at"            example:"2024-06-28"`
	TotalClicks   int                  `json:"total_clicks"          example:"42"`
	UniqueIPs     int                  `json:"unique_ips"            example:"15"`
	LastClickedAt *http_types.DateOnly `json:"last_clicked_at,omitempty" example:"2024-06-28"`
}

func statsDTOFromDomain(stats domain.LinkStats) StatsDTOResponse {
	var lastClickedAt *http_types.DateOnly
	if stats.LastClickedAt != nil {
		d := http_types.DateOnly(*stats.LastClickedAt)
		lastClickedAt = &d
	}

	return StatsDTOResponse{
		ShortURL:      stats.ShortURL,
		OriginalURL:   stats.OriginalURL,
		CreatedAt:     http_types.DateOnly(stats.CreatedAt),
		TotalClicks:   stats.TotalClicks,
		UniqueIPs:     stats.UniqueIPs,
		LastClickedAt: lastClickedAt,
	}
}

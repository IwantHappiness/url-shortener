package stats_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	http_server "github.com/IwantHappiness/url-shortener/internal/core/transport/http/server"
)

type StatsHTTPHandler struct {
	statsService StatsService
}

type StatsService interface {
	GetStats(ctx context.Context, shortURL string, from, to *time.Time) (domain.LinkStats, error)
}

func NewStatsHTTPHandler(statsService StatsService) *StatsHTTPHandler {
	return &StatsHTTPHandler{statsService: statsService}
}

func (h *StatsHTTPHandler) Routes() []http_server.Route {
	return []http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/stats",
			Handler: h.GetStats,
		},
	}
}

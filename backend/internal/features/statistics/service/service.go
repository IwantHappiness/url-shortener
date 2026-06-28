package stats_service

import (
	"context"
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type StatsService struct {
	statsRepository StatsRepository
}

type StatsRepository interface {
	GetStats(ctx context.Context, shortURL string, from, to *time.Time) (domain.LinkStats, error)
}

func NewStatsService(statsRepository StatsRepository) *StatsService {
	return &StatsService{
		statsRepository: statsRepository,
	}
}

package stats_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
)

func (s *StatsService) GetStats(ctx context.Context, shortURL string, from, to *time.Time) (domain.LinkStats, error) {
	shortURL = strings.TrimSpace(shortURL)
	if shortURL == "" {
		return domain.LinkStats{}, fmt.Errorf("shortURL must not be empty: %w", core_errors.ErrInvalidArgument)
	}

	if from != nil && to != nil && from.After(*to) {
		return domain.LinkStats{}, fmt.Errorf("'from' must be before or equal to 'to': %w", core_errors.ErrInvalidArgument)
	}

	stats, err := s.statsRepository.GetStats(ctx, shortURL, from, to)
	if err != nil {
		return domain.LinkStats{}, fmt.Errorf("get stats from repository: %w", err)
	}

	return stats, nil
}

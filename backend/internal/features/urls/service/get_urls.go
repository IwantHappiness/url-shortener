package url_service

import (
	"context"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
)

func (s *URLService) GetURLs(ctx context.Context, userID, limit, offset *int) ([]domain.Link, error) {
	if userID != nil && *userID < 0 {
		return nil, fmt.Errorf("userID must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	urls, err := s.urlRepository.GetURLs(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get urls from repository: %w", err)
	}

	return urls, nil
}

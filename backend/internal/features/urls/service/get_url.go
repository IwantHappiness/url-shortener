package url_service

import (
	"context"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

func (s *URLService) GetURL(ctx context.Context, urlID int) (domain.Link, error) {
	link, err := s.urlRepository.GetURL(ctx, urlID)
	if err != nil {
		return domain.Link{}, fmt.Errorf("failed to get url from repository: %w", err)
	}

	return link, nil
}

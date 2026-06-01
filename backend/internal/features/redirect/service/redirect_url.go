package redirect_service

import (
	"context"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

func (s *RedirectService) RedirectURL(ctx context.Context, shortURL string) (domain.Link, error) {
	link, err := s.redirectRepository.GetByShortURL(ctx, shortURL)
	if err != nil {
		return domain.Link{}, fmt.Errorf("get original ulr: %w", err)
	}
	return link, nil
}

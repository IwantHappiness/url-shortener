package url_service

import (
	"context"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (s *URLService) PatchURL(ctx context.Context, urlID int) (domain.Link, error) {
	link, err := s.urlRepository.GetURL(ctx, urlID)
	if err != nil {
		return domain.Link{}, fmt.Errorf("get url: %w", err)
	}

	randomURL, err := gonanoid.New(12)
	if err != nil {
		return domain.Link{}, fmt.Errorf("generate short url: %w", err)
	}

	link.ApplyPatch(randomURL)

	patchedLink, err := s.urlRepository.PatchURL(ctx, urlID, link)
	if err != nil {
		return domain.Link{}, fmt.Errorf("patch url: %w", err)
	}

	return patchedLink, nil
}

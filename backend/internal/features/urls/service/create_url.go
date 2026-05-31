package url_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (u *URLService) CreateURL(ctx context.Context, link domain.Link) (domain.Link, error) {
	if err := link.ValidateLink(); err != nil {
		return domain.Link{}, fmt.Errorf("validate link domain: %w", err)
	}

	var err error

	for range MaxAttempts {
		randomURL, genErr := gonanoid.New(12)
		if genErr != nil {
			return domain.Link{}, fmt.Errorf("generate short URL: %w", err)
		}

		link.ShortURL = randomURL

		link, err = u.urlRepository.CreateURL(ctx, link)
		if err == nil {
			return link, nil
		}

		if errors.Is(err, core_postgres_pool.ErrUniqueViolation) {
			continue
		}

		return domain.Link{}, fmt.Errorf("create short URL: %w", err)
	}

	return domain.Link{}, fmt.Errorf("failed to generate unique short url after %d attempts: %w", MaxAttempts, err)
}

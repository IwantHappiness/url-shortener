package redirect_service

import (
	"context"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type RedirectService struct {
	redirectRepository RedirectRepository
}

type RedirectRepository interface {
	GetByShortURL(ctx context.Context, shortURL string) (domain.Link, error)
}

func NewRedirectService(redirectRepository RedirectRepository) *RedirectService {
	return &RedirectService{
		redirectRepository: redirectRepository,
	}
}

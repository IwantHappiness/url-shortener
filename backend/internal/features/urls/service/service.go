package url_service

import (
	"context"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

const MaxAttempts int = 6

type URLService struct {
	urlRepository URLRepository
}

type URLRepository interface {
	CreateURL(ctx context.Context, link domain.Link) (domain.Link, error)
	GetURLs(ctx context.Context, userID, limit, offset *int) ([]domain.Link, error)
	GetURL(ctx context.Context, id int) (domain.Link, error)
	DeleteURL(ctx context.Context, id int) error
	PatchURL(ctx context.Context, id int, link domain.Link) (domain.Link, error)
}

func NewURLService(urlRepository URLRepository) *URLService {
	return &URLService{
		urlRepository: urlRepository,
	}
}

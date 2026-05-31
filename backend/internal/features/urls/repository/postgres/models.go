package url_postgres_repository

import (
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type LinkModel struct {
	ID          int
	UserID      int
	Version     int
	OriginalURL string
	ShortURL    string
	CreatedAt   time.Time
}

func urlDomainFromModel(model LinkModel) domain.Link {
	return domain.Link{
		ID:          model.ID,
		UserID:      model.UserID,
		Version:     model.Version,
		OriginalURL: model.OriginalURL,
		ShortURL:    model.ShortURL,
		CreatedAt:   model.CreatedAt,
	}
}

func ulrsDomainsFromModels(models []LinkModel) []domain.Link {
	urls := make([]domain.Link, len(models))

	for i, model := range models {
		urls[i] = urlDomainFromModel(model)
	}

	return urls
}

package redirect_postgres_repository

import "github.com/IwantHappiness/url-shortener/internal/core/domain"

type RedirectModel struct {
	OriginalURL string
	ShortURL    string
}

func redirectDomainFromModel(model RedirectModel) domain.Link {
	return domain.Link{
		OriginalURL: model.OriginalURL,
		ShortURL:    model.ShortURL,
	}
}

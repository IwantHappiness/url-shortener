package urls_transport_http

import (
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type UrlDTOResponse struct {
	ID          int       `json:"id"`
	Version     int       `json:"version"`
	UserID      int       `json:"user_id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func urlDTOfromDomain(link domain.Link) UrlDTOResponse {
	return UrlDTOResponse{
		ID:          link.ID,
		Version:     link.Version,
		UserID:      link.UserID,
		OriginalURL: link.OriginalURL,
		ShortURL:    link.ShortURL,
		CreatedAt:   link.CreatedAt,
	}
}

func urlsDTOFromDomains(ulrs []domain.Link) []UrlDTOResponse {
	usersDTO := make([]UrlDTOResponse, len(ulrs))

	for i, user := range ulrs {
		usersDTO[i] = urlDTOfromDomain(user)
	}

	return usersDTO
}

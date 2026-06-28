package urls_transport_http

import (
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type UrlDTOResponse struct {
	ID          int       `json:"id"           example:"1"`
	Version     int       `json:"version"      example:"1"`
	UserID      int       `json:"user_id"      example:"1"`
	OriginalURL string    `json:"original_url" example:"https://example.com"`
	ShortURL    string    `json:"short_url"    example:"abc123"`
	CreatedAt   time.Time `json:"created_at"   example:"2024-06-28T15:04:05Z07:00"`
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

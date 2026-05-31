package domain

import (
	"fmt"
	"net/url"
	"time"

	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
)

type Link struct {
	ID          int
	UserID      int
	Version     int
	OriginalURL string
	ShortURL    string
	CreatedAt   time.Time
}

func NewLink(id, version, authorUserID int, originalURL, shortURL string, createdAt time.Time) Link {
	return Link{
		ID:          id,
		Version:     version,
		UserID:      authorUserID,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		CreatedAt:   createdAt,
	}
}

func NewUninitializedLink(authorUserID int, originalURL string) Link {
	return NewLink(UninitializedID, UninitializedVersion, authorUserID, originalURL, UninitializedShortURL, UninitializedCreatedAt)
}

func (l *Link) ApplyPatch(shortURL string) {
	tmp := *l

	tmp.ShortURL = shortURL

	*l = tmp
}

func (l *Link) ValidateLink() error {
	if l.UserID < 0 {
		return fmt.Errorf("invalid author user id: %w", core_errors.ErrInvalidArgument)
	}

	if l.OriginalURL == "" {
		return fmt.Errorf("original url is empty: %w", core_errors.ErrInvalidArgument)
	}

	if len([]rune(l.OriginalURL)) >= 2048 {
		return fmt.Errorf("original url too long: %w", core_errors.ErrInvalidArgument)
	}

	parsedURL, err := url.Parse(l.OriginalURL)
	if err != nil {
		return fmt.Errorf("parse original url: %w", core_errors.ErrInvalidArgument)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid url scheme: %w", core_errors.ErrInvalidArgument)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("url host is empty: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

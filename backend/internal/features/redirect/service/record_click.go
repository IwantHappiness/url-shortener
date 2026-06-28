package redirect_service

import (
	"context"
	"fmt"
	"net"
	"strings"

	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
)

func (s *RedirectService) RecordClick(ctx context.Context, shortURL, ip string) error {
	shortURL = strings.TrimSpace(shortURL)
	if shortURL == "" {
		return fmt.Errorf("shortURL must not be empty: %w", core_errors.ErrInvalidArgument)
	}

	ip = strings.TrimSpace(ip)
	if ip == "" {
		return fmt.Errorf("client IP must not be empty: %w", core_errors.ErrInvalidArgument)
	}

	if parsedIP := net.ParseIP(ip); parsedIP == nil {
		return fmt.Errorf("invalid client IP='%s': %w", ip, core_errors.ErrInvalidArgument)
	}

	if err := s.redirectRepository.RecordClick(ctx, shortURL, ip); err != nil {
		return fmt.Errorf("record click in repository: %w", err)
	}

	return nil
}

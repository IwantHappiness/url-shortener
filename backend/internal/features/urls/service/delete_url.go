package url_service

import (
	"context"
	"fmt"
)

func (s *URLService) DeleteURL(ctx context.Context, urlID int) error {
	if err := s.urlRepository.DeleteURL(ctx, urlID); err != nil {
		return fmt.Errorf("delete url form repository: %w", err)
	}

	return nil
}

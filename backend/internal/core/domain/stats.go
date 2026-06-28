package domain

import "time"

type LinkStats struct {
	ShortURL      string
	OriginalURL   string
	CreatedAt     time.Time
	TotalClicks   int
	UniqueIPs     int
	LastClickedAt *time.Time
}

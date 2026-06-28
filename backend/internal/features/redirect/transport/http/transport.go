package redirect_transport_http

import (
	"context"
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type RedirectHTTPHandler struct {
	redirectService RedirectService
}

type RedirectService interface {
	RedirectURL(ctx context.Context, shortCode string) (domain.Link, error)
	RecordClick(ctx context.Context, shortCode, ip string) error
}

func NewRedirectHTTPHandler(redirectService RedirectService) *RedirectHTTPHandler {
	return &RedirectHTTPHandler{
		redirectService: redirectService,
	}
}

func (h *RedirectHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.RedirectURL(w, r)
}

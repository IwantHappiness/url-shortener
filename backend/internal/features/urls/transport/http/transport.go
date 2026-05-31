package urls_transport_http

import (
	"context"
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	http_server "github.com/IwantHappiness/url-shortener/internal/core/transport/http/server"
)

type UrlsHTTPHandler struct {
	urlsService UrlsService
}

type UrlsService interface {
	CreateURL(ctx context.Context, link domain.Link) (domain.Link, error)
	GetURLs(ctx context.Context, userID, limit, offset *int) ([]domain.Link, error)
	GetURL(ctx context.Context, ulrID int) (domain.Link, error)
	DeleteURL(ctx context.Context, urlID int) error
	PatchURL(ctx context.Context, urlID int) (domain.Link, error)
}

func NewUrlsHTTPHandler(urlsService UrlsService) *UrlsHTTPHandler {
	return &UrlsHTTPHandler{
		urlsService: urlsService,
	}
}

func (h *UrlsHTTPHandler) Routes() []http_server.Route {
	return []http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/urls",
			Handler: h.CreateURL,
		},
		{
			Method:  http.MethodGet,
			Path:    "/urls",
			Handler: h.GetURLs,
		},
		{
			Method:  http.MethodGet,
			Path:    "/urls/{id}",
			Handler: h.GetURL,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/urls/{id}",
			Handler: h.DeleteURL,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/urls/{id}",
			Handler: h.PatchURL,
		},
	}
}

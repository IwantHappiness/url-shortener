package http_server

import (
	"net/http"

	http_middleware "github.com/IwantHappiness/url-shortener/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []http_middleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return http_middleware.ChainMiddleware(r.Handler, r.Middleware...)
}

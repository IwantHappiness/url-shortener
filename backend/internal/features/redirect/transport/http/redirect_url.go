package redirect_transport_http

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
	"go.uber.org/zap"
)

func (rh *RedirectHTTPHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	shortURL := r.PathValue("shortURL")
	if shortURL == "" {
		responseHandler.ErrorResponse(
			fmt.Errorf("short ulr is empty: %w", core_errors.ErrInvalidArgument),
			"failed to get url path value",
		)
		return
	}

	link, err := rh.redirectService.RedirectURL(ctx, shortURL)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to resolve short url",
		)
		return
	}

	clientIP := getClientIP(r)
	if clientIP == "" {
		log.Warn("failed to detect client ip for click statistics", zap.String("shortURL", shortURL))
	} else if err := rh.redirectService.RecordClick(ctx, shortURL, clientIP); err != nil {
		log.Warn(
			"failed to record redirect click",
			zap.String("shortURL", shortURL),
			zap.String("clientIP", clientIP),
			zap.Error(err),
		)
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
}

func getClientIP(r *http.Request) string {
	if forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwardedFor != "" {
		if firstIP, _, _ := strings.Cut(forwardedFor, ","); firstIP != "" {
			return strings.TrimSpace(firstIP)
		}
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	if remoteAddr == "" {
		return ""
	}

	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return host
	}

	return remoteAddr
}

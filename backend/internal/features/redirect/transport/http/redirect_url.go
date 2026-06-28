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

// RedirectURL		godoc
// @Summary 		Редирект по короткой ссылке
// @Description 	Перенаправить пользователя на оригинальный URL по короткому коду
// @Tags 			redirect
// @Accept 			json
// @Produce 		json
// @Param 			shortURL path string true "Short URL code"
// @Success 		307 "Редирект на оригинальный URL"
// @Failure 		400 {object} http_response.ErrorResponse "Bad request"
// @Failure 		404 {object} http_response.ErrorResponse "Not found"
// @Failure 		500 {object} http_response.ErrorResponse "Internal server error"
// @Router 			/{shortURL} [get]
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

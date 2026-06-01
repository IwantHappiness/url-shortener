package redirect_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
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

	http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
}

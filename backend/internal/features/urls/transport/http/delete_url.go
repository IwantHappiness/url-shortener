package urls_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

func (h *UrlsHTTPHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	urlID, err := http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get urlID path value")
		return
	}

	if err = h.urlsService.DeleteURL(ctx, urlID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete url")
		return
	}

	responseHandler.NoContentResponse()
}

package urls_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type GetURLResponse UrlDTOResponse

func (h *UrlsHTTPHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	urlID, err := http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get urlID path value")
		return
	}

	urlDomain, err := h.urlsService.GetURL(ctx, urlID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get url")
		return
	}

	response := GetURLResponse(urlDTOfromDomain(urlDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

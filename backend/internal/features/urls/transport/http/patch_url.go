package urls_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type PatchLinkRequest struct{}

type PatchLinkResponse UrlDTOResponse

func (h *UrlsHTTPHandler) PatchURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	urlID, err := http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get urlID path value")
		return
	}

	linkDomain, err := h.urlsService.PatchURL(ctx, urlID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch url")
		return
	}

	response := PatchLinkResponse(urlDTOfromDomain(linkDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

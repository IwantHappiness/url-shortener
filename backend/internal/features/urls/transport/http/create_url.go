package urls_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type CreateUrlRequest struct {
	URL    string `json:"url" validate:"required,url"`
	UserID int    `json:"user_id" validate:"required"`
}

type CreateUrlResponse UrlDTOResponse

func (h *UrlsHTTPHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke CreateUrl handler")

	var req CreateUrlRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	linkDomain := domainFromDTO(req)

	linkDomain, err := h.urlsService.CreateURL(ctx, linkDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create short URL")
		return
	}

	response := CreateUrlResponse(urlDTOfromDomain(linkDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUrlRequest) domain.Link {
	return domain.NewUninitializedLink(dto.UserID, dto.URL)
}

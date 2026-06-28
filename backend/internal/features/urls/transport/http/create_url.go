package urls_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type CreateUrlRequest struct {
	URL    string `json:"url"      validate:"required,url" example:"https://example.com"`
	UserID int    `json:"user_id"  validate:"required"     example:"1"`
}

type CreateUrlResponse UrlDTOResponse

// CreateURL		godoc
// @Summary 		Создать короткую ссылку
// @Description 	Создать новую короткую ссылку для указанного URL
// @Tags 			urls
// @Accept 			json
// @Produce 		json
// @Param 			request body CreateUrlRequest true "CreateUrl тело запроса"
// @Success 		201 {object} CreateUrlResponse "Успешно созданная ссылка"
// @Failure 		400 {object} http_response.ErrorResponse "Bad request"
// @Failure 		409 {object} http_response.ErrorResponse "Conflict"
// @Failure 		500 {object} http_response.ErrorResponse "Internal server error"
// @Router 			/urls [post]
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

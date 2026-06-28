package urls_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type PatchLinkRequest struct{}

type PatchLinkResponse UrlDTOResponse

// PatchURL		godoc
// @Summary 		Обновить короткую ссылку
// @Description 	Сгенерировать новый короткий код для существующей ссылки
// @Tags 			urls
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "URL ID"
// @Success 		200 {object} PatchLinkResponse "Обновлённая ссылка"
// @Failure 		400 {object} http_response.ErrorResponse "Bad request"
// @Failure 		404 {object} http_response.ErrorResponse "Not found"
// @Failure 		500 {object} http_response.ErrorResponse "Internal server error"
// @Router 			/urls/{id} [patch]
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

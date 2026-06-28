package urls_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

// DeleteURL		godoc
// @Summary 		Удалить ссылку
// @Description 	Удалить короткую ссылку по ID
// @Tags 			urls
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "URL ID"
// @Success 		204 "Ссылка удалена"
// @Failure 		400 {object} http_response.ErrorResponse "Bad request"
// @Failure 		404 {object} http_response.ErrorResponse "Not found"
// @Failure 		500 {object} http_response.ErrorResponse "Internal server error"
// @Router 			/urls/{id} [delete]
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

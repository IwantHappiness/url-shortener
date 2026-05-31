package urls_transport_http

import (
	"fmt"
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type GetURLsResponse []UrlDTOResponse

func (h *UrlsHTTPHandler) GetURLs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userID, limit, offset, err := GetUserIdLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id/limit/offset query params")
		return
	}

	urlsDomain, err := h.urlsService.GetURLs(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get urls")
		return
	}

	response := GetURLsResponse(urlsDTOFromDomains(urlsDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func GetUserIdLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := http_request.GetIntQueryParams(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err := http_request.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}

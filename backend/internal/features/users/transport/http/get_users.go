package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	limit, offset, err := GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit/offset query params")
		return
	}

	users, err := h.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	usersDTO := GetUsersResponse(usersDTOFromDomains(users))

	responseHandler.JSONResponse(usersDTO, http.StatusOK)
}

func GetLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, error := http_request.GetIntQueryParams(r, "limit")
	if error != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", error)
	}

	offset, error := http_request.GetIntQueryParams(r, "offset")
	if error != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", error)
	}

	return limit, offset, nil
}

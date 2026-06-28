package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers		godoc
// @Summary 		Получить список пользователей
// @Description 	Получить список всех пользователей с пагинацией
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			limit query int false "limit"
// @Param 			offset query int false "offset"
// @Success 		200 {array} UserDTOResponse "Список пользователей"
// @Failure 		400 {object} http_response.ErrorResponse "Bad request"
// @Failure 		500 {object} http_response.ErrorResponse "Internal server error"
// @Router 			/users [get]
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
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := http_request.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}

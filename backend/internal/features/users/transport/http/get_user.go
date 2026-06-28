package users_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

// GetUser		godoc
// @Summary 		Получить пользователя по ID
// @Description 	Получить информацию о конкретном пользователе
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "User ID"
// @Success 		200 {object} GetUserResponse "Информация о пользователе"
// @Failure 		400 {object} http_response.ErrorResponse "Bad request"
// @Failure 		404 {object} http_response.ErrorResponse "Not found"
// @Failure 		500 {object} http_response.ErrorResponse "Internal server error"
// @Router 			/users/{id} [get]
func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userId, err := http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	user, err := h.userService.GetUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}

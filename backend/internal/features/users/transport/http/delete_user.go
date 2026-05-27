package users_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
	http_utils "github.com/IwantHappiness/url-shortener/internal/core/transport/http/utils"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userId, err := http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user id path value")
		return
	}

	if err := h.userService.DeleteUser(ctx, userId); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	responseHandler.NoContentResponse()
}

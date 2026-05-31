package users_transport_http

import (
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	Nickname string `json:"nickname" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
}

type CreateUserResponse UserDTOResponse

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke CreateUser handler")

	var req CreateUserRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	userDomain := domainFromDTO(req)

	userDomain, err := h.userService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.Nickname, dto.Email)
}

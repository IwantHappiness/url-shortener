package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_request "github.com/IwantHappiness/url-shortener/internal/core/transport/http/request"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
	http_types "github.com/IwantHappiness/url-shortener/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	Nickname http_types.Nullable[string] `json:"nickname"`
	Email    http_types.Nullable[string] `json:"email"`
}

func (r *PatchUserRequest) Validate() error {
	if r.Nickname.Set {
		if r.Nickname.Value == nil {
			return fmt.Errorf("'Nickname' can't be NULL")
		}

		nicknameLen := len([]rune(*r.Nickname.Value))

		if nicknameLen < 3 || nicknameLen > 20 {
			return fmt.Errorf("'Nickname' must be between 3 and 20 characters")
		}
	}

	if r.Email.Set {
		if r.Email.Value == nil {
			return fmt.Errorf("'Email' can't be NULL")
		}

		emailLen := len([]rune(*r.Email.Value))

		if emailLen < 3 || emailLen > 20 {
			return fmt.Errorf("'Email' must be between 3 and 254 characters")
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser		godoc
// @Summary 		Частично обновить пользователя
// @Description 	Обновить nickname и/или email пользователя. Отсутствующие поля не меняются.
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "User ID"
// @Param 			request body PatchUserRequest true "Patch тело запроса"
// @Success 		200 {object} PatchUserResponse "Обновлённый пользователь"
// @Failure 		400 {object} http_response.ErrorResponse "Bad request"
// @Failure 		404 {object} http_response.ErrorResponse "Not found"
// @Failure 		409 {object} http_response.ErrorResponse "Conflict"
// @Failure 		500 {object} http_response.ErrorResponse "Internal server error"
// @Router 			/users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	userId, err := http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user id path value")
		return
	}

	var req PatchUserRequest
	if err := http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(req)

	userDomain, err := h.userService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(req PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(req.Nickname.ToDomain(), req.Email.ToDomain())
}

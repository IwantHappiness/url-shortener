package users_transport_http

import (
	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	http_types "github.com/IwantHappiness/url-shortener/internal/core/transport/http/types"
)

type UserDTOResponse struct {
	ID        int                 `json:"id"          example:"1"`
	Version   int                 `json:"version"     example:"1"`
	Nickname  string              `json:"nickname"    example:"Иван Иванов"`
	Email     string              `json:"email"       example:"ivanov@gmail.com"`
	CreatedAt http_types.DateOnly `json:"created_at"  example:"2024-06-28"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:        user.ID,
		Version:   user.Version,
		Nickname:  user.Nickname,
		Email:     user.Email,
		CreatedAt: http_types.DateOnly(user.CreatedAt),
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}

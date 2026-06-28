package users_transport_http

import (
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type UserDTOResponse struct {
	ID        int       `json:"id"          example:"1"`
	Version   int       `json:"version"     example:"1"`
	Nickname  string    `json:"nickname"    example:"Иван Иванов"`
	Email     string    `json:"email"       example:"ivanov@gmail.com"`
	CreatedAt time.Time `json:"created_at"  example:"2024-06-28T15:04:05Z07:00"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:        user.ID,
		Version:   user.Version,
		Nickname:  user.Nickname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}

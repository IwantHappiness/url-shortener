package users_transport_http

import (
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type UserDTOResponse struct {
	ID        int       `json:"id"`
	Version   int       `json:"version"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
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

package users_postgres_repository

import (
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type UserModel struct {
	ID        int       `json:"id"`
	Version   int       `json:"version"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func userDomainsFromModels(userModels []UserModel) []domain.User {
	users := make([]domain.User, len(userModels))

	for i, userModel := range userModels {
		users[i] = domain.User{
			ID:        userModel.ID,
			Version:   userModel.Version,
			Nickname:  userModel.Nickname,
			Email:     userModel.Email,
			CreatedAt: userModel.CreatedAt,
		}
	}

	return users
}

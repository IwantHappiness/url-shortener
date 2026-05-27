package user_service

import (
	"context"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

type UserService struct {
	userRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error)
}

func NewUserService(userRepository UsersRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

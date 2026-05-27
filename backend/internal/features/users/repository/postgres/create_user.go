package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `INSERT INTO url_shortener.users (nickname, email) VALUES ($1, $2) RETURNING id, version, nickname, email, created_at;`

	row := r.pool.QueryRow(ctx, query, user.Nickname, user.Email)

	var userModel UserModel
	err := row.Scan(&userModel.ID, &userModel.Version, &userModel.Nickname, &userModel.Email, &userModel.CreatedAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Nickname,
		userModel.Email,
		userModel.CreatedAt,
	), nil
}

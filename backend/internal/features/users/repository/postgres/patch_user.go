package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
)

func (r *UserRepository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE url_shortener.users
	SET nickname = $1, email = $2, version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING id, version, nickname, email, created_at;`

	row := r.pool.QueryRow(ctx, query, user.Nickname, user.Email, user.ID, user.Version)

	var userModel UserModel
	err := row.Scan(&userModel.ID, &userModel.Version, &userModel.Nickname, &userModel.Email, &userModel.CreatedAt)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d' concurrently accessed: %w", id, core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(userModel.ID, userModel.Version, userModel.Nickname, userModel.Email, userModel.CreatedAt)

	return userDomain, nil
}

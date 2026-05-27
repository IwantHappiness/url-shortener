package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UserRepository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE url_shortener.users
	SET nickname = $1, email = $2, version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING id, version, nickname, email, created_at`

	row := r.pool.QueryRow(ctx, query, user.Nickname, user.Email, user.ID, user.Version)

	var UserModel domain.User
	err := row.Scan(&UserModel.ID, &UserModel.Version, &UserModel.Nickname, &UserModel.Email, &UserModel.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d' concurrently accessed: %w", id, core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(UserModel.ID, UserModel.Version, UserModel.Nickname, UserModel.Email, UserModel.CreatedAt)

	return userDomain, nil
}

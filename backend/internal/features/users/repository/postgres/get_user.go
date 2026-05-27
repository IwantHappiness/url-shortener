package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UserRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT id, version, nickname, email, created_at FROM url_shortener.users WHERE id = $1;`

	row := r.pool.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Version, &user.Nickname, &user.Email, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		user.ID,
		user.Version,
		user.Nickname,
		user.Email,
		user.CreatedAt,
	)

	return userDomain, nil
}

package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/IwantHappiness/url-shortener/internal/core/domain"
)

func (r *UserRepository) GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT id, version, nickname, email, created_at FROM url_shortener.users ORDER BY id ASC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var userModels []UserModel

	for rows.Next() {
		var userModel UserModel
		if err := rows.Scan(&userModel.ID, &userModel.Version, &userModel.Nickname, &userModel.Email, &userModel.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		userModels = append(userModels, userModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}

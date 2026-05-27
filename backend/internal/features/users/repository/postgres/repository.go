package users_postgres_repository

import core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"

type UserRepository struct {
	pool core_postgres_pool.ConnectionPool
}

func NewUserRepository(pool *core_postgres_pool.ConnectionPool) *UserRepository {
	return &UserRepository{pool: *pool}
}

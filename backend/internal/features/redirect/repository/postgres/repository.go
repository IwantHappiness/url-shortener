package redirect_postgres_repository

import core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"

type RedirectRepository struct {
	pool core_postgres_pool.Pool
}

func NewRedirectRepository(pool core_postgres_pool.Pool) *RedirectRepository {
	return &RedirectRepository{pool: pool}
}

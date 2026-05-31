package url_postgres_repository

import core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"

type URLRepository struct {
	pool core_postgres_pool.Pool
}

func NewURLRepository(pool core_postgres_pool.Pool) *URLRepository {
	return &URLRepository{pool: pool}
}

package stats_postgres_repository

import core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"

type StatsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatsRepository(pool core_postgres_pool.Pool) *StatsRepository {
	return &StatsRepository{pool: pool}
}

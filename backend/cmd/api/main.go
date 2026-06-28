package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/IwantHappiness/url-shortener/internal/core/config"
	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	core_pgx_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool/pgx"
	http_middleware "github.com/IwantHappiness/url-shortener/internal/core/transport/http/middleware"
	http_server "github.com/IwantHappiness/url-shortener/internal/core/transport/http/server"
	redirect_postgres_repository "github.com/IwantHappiness/url-shortener/internal/features/redirect/repository/postgres"
	redirect_service "github.com/IwantHappiness/url-shortener/internal/features/redirect/service"
	redirect_transport_http "github.com/IwantHappiness/url-shortener/internal/features/redirect/transport/http"
	stats_postgres_repository "github.com/IwantHappiness/url-shortener/internal/features/statistics/repository/postgres"
	stats_service "github.com/IwantHappiness/url-shortener/internal/features/statistics/service"
	stats_transport_http "github.com/IwantHappiness/url-shortener/internal/features/statistics/transport/http"
	url_postgres_repository "github.com/IwantHappiness/url-shortener/internal/features/urls/repository/postgres"
	url_service "github.com/IwantHappiness/url-shortener/internal/features/urls/service"
	urls_transport_http "github.com/IwantHappiness/url-shortener/internal/features/urls/transport/http"
	users_postgres_repository "github.com/IwantHappiness/url-shortener/internal/features/users/repository/postgres"
	user_service "github.com/IwantHappiness/url-shortener/internal/features/users/service"
	users_transport_http "github.com/IwantHappiness/url-shortener/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Debug("application time zone", zap.Any("timeZone", time.Local))

	log.Debug("initialized postgres connection pool")

	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		log.Fatal("failed to init Postgres connection pool:", zap.Error(err))
	}
	defer pool.Close()

	log.Debug("initialized feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUserRepository(pool)
	usersService := user_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(usersService)

	log.Debug("initialized feature", zap.String("feature", "urls"))

	urlsRepository := url_postgres_repository.NewURLRepository(pool)
	urlsService := url_service.NewURLService(urlsRepository)
	urlsTransportHTTP := urls_transport_http.NewUrlsHTTPHandler(urlsService)

	log.Debug("initilized feature", zap.String("feature", "redirect"))

	redirectRepository := redirect_postgres_repository.NewRedirectRepository(pool)
	redirectService := redirect_service.NewRedirectService(redirectRepository)
	redirectTransportHTTP := redirect_transport_http.NewRedirectHTTPHandler(redirectService)

	log.Debug("initialized feature", zap.String("feature", "statistics"))

	statsRepository := stats_postgres_repository.NewStatsRepository(pool)
	statsService := stats_service.NewStatsService(statsRepository)
	statsTransportHTTP := stats_transport_http.NewStatsHTTPHandler(statsService)

	log.Debug("initialized HTTP server")

	httpServer := http_server.NewHTTPServer(
		http_server.NewConfigMust(), log,
		http_middleware.RequestId(),
		http_middleware.Logger(log),
		http_middleware.Trace(),
		http_middleware.Panic(),
	)

	httpServer.RegisterRedirectHandler("/{shortURL}", redirectTransportHTTP)

	apiVersionRouterV1 := http_server.NewAPIVersionRouter(http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoute(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoute(urlsTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoute(statsTransportHTTP.Routes()...)

	// Example of usage apVersionRouterV2 witch separate Middlewares
	//
	// apiVersionRouterV2 := http_server.NewAPIVersionRouter(
	// 	http_server.ApiVersion2,
	// 	http_middleware.Dummy("api v2 middleware"),
	// )
	// apiVersionRouterV2.RegisterRoute(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouter(
		apiVersionRouterV1,
		// apiVersionRouterV2
	)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}

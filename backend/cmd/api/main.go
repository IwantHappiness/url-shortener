package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	core_postgres_pool "github.com/IwantHappiness/url-shortener/internal/core/repository/postgres/pool"
	http_middleware "github.com/IwantHappiness/url-shortener/internal/core/transport/http/middleware"
	http_server "github.com/IwantHappiness/url-shortener/internal/core/transport/http/server"
	users_postgres_repository "github.com/IwantHappiness/url-shortener/internal/features/users/repository/postgres"
	user_service "github.com/IwantHappiness/url-shortener/internal/features/users/service"
	users_transport_http "github.com/IwantHappiness/url-shortener/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Debug("initialized postgres connection pool")

	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		log.Fatal("failed to init Postgres connection pool:", zap.Error(err))
	}
	defer pool.Close()

	// log.Debug("Starting URL-Shortener application")

	log.Debug("initialized feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUserRepository(pool)
	usersService := user_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(usersService)

	log.Debug("initialized HTTP server")

	httpServer := http_server.NewHTTPServer(
		http_server.NewConfigMust(), log,
		http_middleware.RequestId(),
		http_middleware.Logger(log),
		http_middleware.Trace(),
		http_middleware.Panic(),
	)

	apiVersionRouterV1 := http_server.NewAPIVersionRouter(http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoute(usersTransportHTTP.Routes()...)

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

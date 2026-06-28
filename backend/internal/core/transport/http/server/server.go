package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/IwantHappiness/url-shortener/docs"
	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_middleware "github.com/IwantHappiness/url-shortener/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *logger.Logger

	middleware []http_middleware.Middleware
}

func NewHTTPServer(config Config, log *logger.Logger, middleware ...http_middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     NewConfigMust(),
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouter(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router.WithMiddleware()))
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DefaultModelsExpandDepth(-1),
	))

	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)

}

func (s *HTTPServer) RegisterRedirectHandler(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr:", s.config.Addr))

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}

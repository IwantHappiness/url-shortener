package http_middleware

import (
	"net/http"
	"time"

	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIdHeader = "X-Request-Id"

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIdHeader)
			if requestId == "" {
				requestId = uuid.NewString()
			}

			r.Header.Set(requestIdHeader, requestId)
			w.Header().Set(requestIdHeader, requestId)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIdHeader)

			l := log.With(
				zap.String("request_id", requestId),
				zap.String("url", r.URL.String()),
			)

			ctx := logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			rw := http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(">>> incoming HTTP request", zap.String("http_method", r.Method), zap.Time("time", before.UTC()))

			next.ServeHTTP(rw, r)

			log.Debug("<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Since(before)))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)

			responseHandler := http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicHandler(p, "during handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

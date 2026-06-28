package stats_transport_http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
	"github.com/IwantHappiness/url-shortener/internal/core/logger"
	http_response "github.com/IwantHappiness/url-shortener/internal/core/transport/http/response"
)

// GetStats		godoc
// @Summary 		Получить статистику по ссылке
// @Description 	Получить статистику переходов по короткой ссылке. Можно указать диапазон дат.
// @Tags 			statistics
// @Accept 			json
// @Produce 		json
// @Param 			short_url query string true "Short URL code"
// @Param 			from query string false "Начало диапазона (RFC3339)"
// @Param 			to query string false "Конец диапазона (RFC3339)"
// @Success 		200 {object} StatsDTOResponse "Статистика по ссылке"
// @Failure 		400 {object} http_response.ErrorResponse "Bad request"
// @Failure 		404 {object} http_response.ErrorResponse "Not found"
// @Failure 		500 {object} http_response.ErrorResponse "Internal server error"
// @Router 			/stats [get]
func (h *StatsHTTPHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := http_response.NewHTTPResponseHandler(log, w)

	shortURL, from, to, err := getShortURLFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get short_url/from/to query params")
		return
	}

	if shortURL == nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("short_url query param is required: %w", core_errors.ErrInvalidArgument),
			"failed to get short_url/from/to query params",
		)
		return
	}

	stats, err := h.statsService.GetStats(ctx, *shortURL, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get stats")
		return
	}

	responseHandler.JSONResponse(statsDTOFromDomain(stats), http.StatusOK)
}

func getShortURLFromToQueryParams(r *http.Request) (*string, *time.Time, *time.Time, error) {
	const (
		shortURLQueryParamKey = "short_url"
		fromQueryParamKey     = "from"
		toQueryParamKey       = "to"
	)

	query := r.URL.Query()

	var shortURL *string
	if rawShortURL := strings.TrimSpace(query.Get(shortURLQueryParamKey)); rawShortURL != "" {
		shortURL = &rawShortURL
	}

	from, err := parseTimeQueryParam(query.Get(fromQueryParamKey), fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, err
	}

	to, err := parseTimeQueryParam(query.Get(toQueryParamKey), toQueryParamKey)
	if err != nil {
		return nil, nil, nil, err
	}

	return shortURL, from, to, nil
}

func parseTimeQueryParam(rawValue, key string) (*time.Time, error) {
	rawValue = strings.TrimSpace(rawValue)
	if rawValue == "" {
		return nil, nil
	}

	parsedAt, err := time.Parse(time.RFC3339, rawValue)
	if err != nil {
		return nil, fmt.Errorf("param='%s' by key='%s' not a valid RFC3339 time: %v: %w", rawValue, key, err, core_errors.ErrInvalidArgument)
	}

	return &parsedAt, nil
}

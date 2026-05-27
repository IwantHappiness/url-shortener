package http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValues := r.PathValue(key)
	if pathValues == "" {
		return 0, fmt.Errorf("no key='%s' in path: %w", key, core_errors.ErrInvalidArgument)
	}

	intValue, err := strconv.Atoi(pathValues)
	if err != nil {
		return 0, fmt.Errorf("path value='%s' by key='%s' no a valid integer: %w", pathValues, key, core_errors.ErrInvalidArgument)
	}

	return intValue, nil
}

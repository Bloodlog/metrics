package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

func validateRequest(r *http.Request) int {
	if r.Method != http.MethodPost {
		return http.StatusBadRequest
	}

	parts := strings.Split(r.RequestURI, "/")
	if len(parts) < 5 {
		return http.StatusNotFound
	}

	if parts[1] != "update" {
		return http.StatusBadRequest
	}
	if parts[2] != "counter" && parts[2] != "gauge" {
		return http.StatusBadRequest
	}

	metricID := parts[4]
	if _, err := strconv.Atoi(metricID); err != nil {
		return http.StatusBadRequest
	}

	return 0
}

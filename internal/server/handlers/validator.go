package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

func validateURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parts := strings.Split(r.RequestURI, "/")
	if len(parts) < 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[1] != "update" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if parts[2] != "counter" && parts[2] != "gauge" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricID := parts[4]
	if _, err := strconv.Atoi(metricID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

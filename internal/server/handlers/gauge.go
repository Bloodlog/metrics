package handlers

import (
	"fmt"
	"metrics/internal/server/repository"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GaugeHandler(memStorage *repository.MemStorage, debug bool) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if status := validateRequest(request); status != 0 {
			response.WriteHeader(status)
			return
		}

		timeStr := time.Now().Format("2006-01-02 15:04:05")
		parts := strings.Split(request.RequestURI, "/")

		metricName := parts[3]
		metricValue, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			if debug {
				log := "[" + timeStr + "] " + request.RequestURI + " " + strconv.Itoa(http.StatusOK)
				fmt.Println(log)
			}
			return
		}

		memStorage.SetGauge(metricName, metricValue)

		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(http.StatusOK)

		if debug {
			log := "[" + timeStr + "] " + request.RequestURI + " " + strconv.Itoa(http.StatusOK)
			fmt.Println(log)
		}
	}
}

package handlers

import (
	"fmt"
	"github.com/Bloodlog/metrics/internal/server/storage"
	"net/http"
	"time"
)

func CounterHandler(memStorage *storage.MemStorage) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		validateURL(response, request)

		memStorage.SetCounter("req", 1)

		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(http.StatusOK)

		fmt.Printf("%s - Status: %d\n", time.Now().Format(time.RFC3339), http.StatusOK)
	}
}

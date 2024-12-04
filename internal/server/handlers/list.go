package handlers

import (
	"fmt"
	"metrics/internal/server/repository"
	"net/http"
)

func ListHandler(memStorage *repository.MemStorage) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		const ErrorText = "failed to write response"
		response.Header().Set("Content-Type", "text/html; charset=utf-8")

		_, err := response.Write([]byte("<html><head><title>Metrics List</title></head><body><h1>Metrics</h1><ul>"))
		if err != nil {
			http.Error(response, ErrorText, http.StatusInternalServerError)
			return
		}

		for name, value := range memStorage.Gauges() {
			_, err := response.Write([]byte(fmt.Sprintf("<li>%s (gauge): %.6f</li>", name, value)))
			if err != nil {
				http.Error(response, ErrorText, http.StatusInternalServerError)
				return
			}
		}

		for name, value := range memStorage.Counters() {
			_, err := response.Write([]byte(fmt.Sprintf("<li>%s (counter): %d</li>", name, value)))
			if err != nil {
				http.Error(response, ErrorText, http.StatusInternalServerError)
				return
			}
		}

		_, err = response.Write([]byte("</ul></body></html>"))
		if err != nil {
			http.Error(response, ErrorText, http.StatusInternalServerError)
			return
		}
	}
}

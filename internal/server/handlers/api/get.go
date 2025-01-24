package api

import (
	"encoding/json"
	"errors"
	"metrics/internal/server/service"
	"net/http"
)

func (h *Handler) GetHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "api GetHandler")
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "application/json")

		var metricGetRequest service.MetricsGetRequest

		if err := json.NewDecoder(request.Body).Decode(&metricGetRequest); err != nil {
			handlerLogger.Infow("Invalid JSON", nameError, err)
			response.WriteHeader(http.StatusBadRequest)

			return
		}

		if metricGetRequest.MType != "counter" && metricGetRequest.MType != "gauge" {
			response.WriteHeader(http.StatusBadRequest)

			return
		}

		metricService := service.NewMetricService(handlerLogger)
		result, err := metricService.Get(ctx, metricGetRequest, h.memStorage)
		if err != nil {
			if errors.Is(err, service.ErrMetricNotFound) {
				handlerLogger.Infoln("Metric not found", err)
				response.WriteHeader(http.StatusNotFound)
				return
			}
			handlerLogger.Infow("error in service", nameError, err)
			response.WriteHeader(http.StatusBadRequest)

			return
		}

		resp, err := json.Marshal(result)
		if err != nil {
			handlerLogger.Infow("error marshal json", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = response.Write(resp)
		if err != nil {
			handlerLogger.Infow("error write response", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

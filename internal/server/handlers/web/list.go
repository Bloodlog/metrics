package web

import (
	"html/template"
	"net/http"
)

// ListHandler .
// @Summary Список метрик
// @Description Генерирует HTML-страницу с перечнем метрик (gauge и counter)
// @Tags Info
// @Produce  text/html
// @Success 200 {string} string "HTML страница с метриками"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router / [get].
func (h *Handler) ListHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "web ListHandler")
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := h.metricService.GetMetrics(ctx)

		tmpl, err := template.New("metrics").Parse(`
			<html>
				<head><title>Metrics List</title></head>
				<body>
					<h1>Metrics</h1>
					<ul>
						{{range $name, $value := .Gauges}}
							<li>{{$name}} (gauge): {{$value}}</li>
						{{end}}
						{{range $name, $value := .Counters}}
							<li>{{$name}} (counter): {{$value}}</li>
						{{end}}
					</ul>
				</body>
			</html>
		`)

		if err != nil {
			handlerLogger.Infow("error parse metrics", "error", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(response, data)
		if err != nil {
			handlerLogger.Infow("error parse metrics", "error", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

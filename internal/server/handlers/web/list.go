package web

import (
	"html/template"
	"net/http"
)

type MetricsData struct {
	Gauges   map[string]float64
	Counters map[string]uint64
}

func (h *Handler) ListHandler() http.HandlerFunc {
	handlerLogger := h.logger.With(nameLogger, "web ListHandler")
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		response.Header().Set("Content-Type", "text/html; charset=utf-8")

		gauges, _ := h.memStorage.Gauges(ctx)
		counters, _ := h.memStorage.Counters(ctx)

		data := MetricsData{
			Gauges:   gauges,
			Counters: counters,
		}

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

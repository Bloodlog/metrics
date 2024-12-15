package handlers

import (
	"html/template"
	"metrics/internal/server/repository"
	"net/http"

	"go.uber.org/zap"
)

type MetricsData struct {
	Gauges   map[string]float64
	Counters map[string]uint64
}

func ListHandler(memStorage *repository.MemStorage, logger zap.SugaredLogger) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/html; charset=utf-8")
		const nameError = "error"

		data := MetricsData{
			Gauges:   memStorage.Gauges(),
			Counters: memStorage.Counters(),
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
			logger.Infow("error parse metrics", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(response, data)
		if err != nil {
			logger.Infow("error parse metrics", nameError, err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

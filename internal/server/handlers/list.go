package handlers

import (
	"html/template"
	"log"
	"metrics/internal/server/repository"
	"net/http"
)

type MetricsData struct {
	Gauges   map[string]float64
	Counters map[string]uint64
}

func ListHandler(memStorage *repository.MemStorage) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		const ErrorText = "failed to write response"
		response.Header().Set("Content-Type", "text/html; charset=utf-8")

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
			log.Printf("error parse metrics: %v", err)
			http.Error(response, ErrorText, http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(response, data)
		if err != nil {
			log.Printf("error parse metrics: %v", err)
			http.Error(response, ErrorText, http.StatusInternalServerError)
			return
		}
	}
}

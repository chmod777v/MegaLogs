package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests",
	}, []string{"method", "status"})

	responseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Response time in seconds",
		Buckets: []float64{0.1, 0.5, 1, 2},
	})
)

type Base struct {
	Number int
	String string
}

func Get(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	req := Base{
		Number: 123,
		String: "Hellow, world",
	}
	bytes, err := json.MarshalIndent(req, "", " ")
	if err != nil {
		slog.Error("Error json.Marshal")
		httpRequests.WithLabelValues("GET", "500").Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
	httpRequests.WithLabelValues("GET", "200").Inc()
	responseTime.Observe(time.Since(start).Seconds())
	slog.Info("Sucessfull post Request")
}
func Post(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req Base
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		httpRequests.WithLabelValues("POST", "500").Inc()
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Erro while receiving data", "ERROR", err.Error())
		return
	}

	httpRequests.WithLabelValues("POST", "200").Inc()
	responseTime.Observe(time.Since(start).Seconds())
	slog.Info("RequestGet", "Data", req)
}
func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Get(w, r)
	case http.MethodPost:
		Post(w, r)
	default:
		httpRequests.WithLabelValues(r.Method, "405").Inc()
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", Handler)

	slog.Info("Server listening", "Host", "localhost:8080")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		slog.Error("Error starting server", "ERROR", err.Error())
	}

}

package metric

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(addres string) error {
	mux := http.NewServeMux() //чтобы на порту крутились только метрики
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addres, mux)
}

var requestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "test_server",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"status"})

func ObserveRequest(d time.Duration, status int) {
	requestMetrics.WithLabelValues(strconv.Itoa(status)).Observe(d.Seconds())
}

/*
Второе значение в каждой паре - это допустимая погрешность:
0.5: 0.05 - медиана с точностью ±5%
0.9: 0.01 - 90-й перцентиль с точностью ±1%
0.99: 0.001 - 99-й перцентиль с точностью ±0.1%

test_server_http_request{quantile="0.5"} 0.023
test_server_http_request{quantile="0.9"} 0.045
test_server_http_request{quantile="0.99"} 0.098
test_server_http_request_sum 12.345
test_server_http_request_count 42

50% запросов обрабатываются быстрее 23ms
90% запросов обрабатываются быстрее 45ms
99% запросов обрабатываются быстрее 98ms
Общее время всех запросов: 12.345 секунд
Общее количество запросов: 42
*/

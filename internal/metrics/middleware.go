package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metricsHandler = promhttp.Handler()

func NewMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				metricsHandler.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			statusCode := 200

			rw := &responseWriter{ResponseWriter: w, statusCode: &statusCode}
			next.ServeHTTP(rw, r)

			duration := time.Since(start).Seconds()
			method := r.Method
			path := r.URL.Path

			HTTPRequestsTotal.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
			HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode *int
}

func (rw *responseWriter) WriteHeader(code int) {
	*rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Init() {
	reg := prometheus.DefaultRegisterer
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())
}

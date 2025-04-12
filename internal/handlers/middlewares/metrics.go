package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/hamillka/avitoTechSpring25/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		start := time.Now()

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path

		metrics.HTTPRequestCount.WithLabelValues(r.Method, path, strconv.Itoa(rec.status)).Inc()
		metrics.HTTPResponseDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

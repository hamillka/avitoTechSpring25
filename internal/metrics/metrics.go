package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Технические метрики
	HTTPRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество HTTP-запросов",
		},
		[]string{"method", "path", "status"},
	)

	HTTPResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Время ответа сервера",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Бизнесовые метрики
	PVZCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pvz_created_total",
			Help: "Количество созданных ПВЗ",
		},
	)

	ReceptionsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "receptions_created_total",
			Help: "Количество созданных приёмок",
		},
	)

	ProductsAdded = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "products_added_total",
			Help: "Количество добавленных товаров",
		},
	)
)

func Register() {
	prometheus.MustRegister(
		HTTPRequestCount,
		HTTPResponseDuration,
		PVZCreated,
		ReceptionsCreated,
		ProductsAdded,
	)
}

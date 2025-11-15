package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests",
			Help: "Total Number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	UrlsShortenedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "urls_shortened_total",
			Help: "total number of urls shortened",
		},
	)

	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "Active_Connections",
			Help: "Number of active connections",
		},
	)

	UrlsResolvedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "urls_resolved_total",
			Help: "total_urls_resolved",
		},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "database_query_duration in seconds",
			Help: "Duration of each DB query in seconds",
		},
		[]string{"query_type"},
	)

	RedisCacheDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "redis_operation_duration_seconds",
			Help: "duratoin of each redis operation in seconds",
		},
		[]string{"operation"},
	)

	CacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "Cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "Cache_miss_total",
			Help: "Total numbe of cache misses",
		},
	)
)

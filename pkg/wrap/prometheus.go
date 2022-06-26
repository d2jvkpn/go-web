package wrap

/*
  https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang
*/
import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheusMonitor(namespace string) gin.HandlerFunc {
	totalRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			// Subsystem: value, // string
			Name: "http_requests_total",
			Help: "Number of get requests",
		},
		[]string{"path"},
	)

	responseStatus := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "response_status",
			Help:      "Status of HTTP response",
		},
		[]string{"status"},
	)

	httpDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_response_time_seconds",
			Help:      "Duration of HTTP requests",
		},
		[]string{"path"},
	)

	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)

	return func(ctx *gin.Context) {
		p := ctx.Request.URL.Path
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(p))

		ctx.Next()

		responseStatus.WithLabelValues(strconv.Itoa(ctx.Writer.Status())).Inc()
		totalRequests.WithLabelValues(p).Inc()

		timer.ObserveDuration()
	}
}

func PrometheusFunc(ctx *gin.Context) {
	handler := promhttp.Handler()
	handler.ServeHTTP(ctx.Writer, ctx.Request)
}

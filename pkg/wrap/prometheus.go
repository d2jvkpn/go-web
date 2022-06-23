package wrap

/*
  https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang
*/
import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheusMonitor(prefix string) gin.HandlerFunc {
	addPrefix := func(str string) string {
		if prefix == "" {
			return str
		}
		return fmt.Sprintf("%s_%s", strings.TrimRight(prefix, "_"), str)
	}

	totalRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: addPrefix("http_requests_total"),
			Help: "Number of get requests",
		},
		[]string{"path"},
	)

	responseStatus := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: addPrefix("response_status"),
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)

	httpDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: addPrefix("http_response_time_seconds"),
			Help: "Duration of HTTP requests",
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

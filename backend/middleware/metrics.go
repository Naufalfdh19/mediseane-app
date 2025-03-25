package middleware

import (
	"montelukast/pkg/metrics"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCounter = promauto.NewCounterVec(

		prometheus.CounterOpts{

			Name: "client_request_count",

			Help: "Total number of requests from client",
		},

		[]string{"method", "route", "status"},
	)
)

func IncrementRequestCount(c *gin.Context) {

	c.Next()
	metrics.RequestCount.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(c.Writer.Status())).Inc()

}

func ObserveMiddleware(c *gin.Context) {
	c.Next()

	reqMethod := c.Request.Method
	route := c.FullPath()
	status := strconv.Itoa(c.Writer.Status())

	labels := prometheus.Labels{
		"method": reqMethod,
		"route":  route,
		"status": status,
	}

	RequestCounter.With(labels).Inc()

}

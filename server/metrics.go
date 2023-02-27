package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var booksCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "books_stored",
	Help: "Current size of books array",
})

var reqByURI = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "gin_uri",
	Help: "gin request by uri",
}, []string{"uri", "method", "code"})

var reqDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "gin_req_dur",
	Help: "req duration",
})

func init() {
	prometheus.MustRegister(booksCount, reqByURI, reqDuration)
}

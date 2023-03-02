package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

const URL = "/metrics"

type Prometheus struct {
	reqCnt *prometheus.CounterVec
	reqDur *prometheus.HistogramVec
	reqSz  prometheus.Summary
	resSz  prometheus.Summary
}

func New() *Prometheus {
	reqCntMetric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "gin",
			Name:      "requests_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method", "handler", "host", "url"},
	)
	reqDurMetric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: "gin",
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
		},
		[]string{"code", "method", "url"},
	)
	reqSz := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: "gin",
			Name:      "request_size_bytes",
			Help:      "The HTTP request sizes in bytes.",
		},
	)
	resSz := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: "gin",
			Name:      "response_size_bytes",
			Help:      "The HTTP response sizes in bytes.",
		},
	)
	if err := prometheus.Register(reqCntMetric); err != nil {
		fmt.Println(err)
	}
	if err := prometheus.Register(reqDurMetric); err != nil {
		fmt.Println(err)
	}
	if err := prometheus.Register(reqSz); err != nil {
		fmt.Println(err)
	}
	if err := prometheus.Register(resSz); err != nil {
		fmt.Println(err)
	}
	return &Prometheus{
		reqCnt: reqCntMetric,
		reqDur: reqDurMetric,
		reqSz:  reqSz,
		resSz:  resSz,
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}

func HandlerFunc(p *Prometheus) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == URL {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		elapsed := float64(time.Since(start)) / float64(time.Second)
		status := strconv.Itoa(c.Writer.Status())
		url := c.Request.URL.Path
		resSz := float64(c.Writer.Size())
		reqSz := computeApproximateRequestSize(c.Request)
		p.reqDur.WithLabelValues(status, c.Request.Method, url).Observe(elapsed)
		p.reqCnt.WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Host, url).Inc()
		p.reqSz.Observe(float64(reqSz))
		p.resSz.Observe(resSz)
	}
}

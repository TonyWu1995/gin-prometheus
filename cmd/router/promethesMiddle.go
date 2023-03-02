package router

import (
	prom "gin-prometheus/pkg/prometheus"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const URL = "/metrics"

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

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}

func HandlerFunc(p *prom.Prometheus) gin.HandlerFunc {
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
		p.ReqDur.WithLabelValues(status, c.Request.Method, url).Observe(elapsed)
		p.ReqCnt.WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Host, url).Inc()
		p.ReqSz.Observe(float64(reqSz))
		p.ResSz.Observe(resSz)
	}
}

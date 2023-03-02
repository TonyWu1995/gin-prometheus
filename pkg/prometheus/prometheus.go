package prometheus

import (
	"gin-prometheus/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
)

type Prometheus struct {
	ReqCnt *prometheus.CounterVec
	ReqDur *prometheus.HistogramVec
	ReqSz  prometheus.Summary
	ResSz  prometheus.Summary
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
		log.Logger.Errorln(err)
	}
	if err := prometheus.Register(reqDurMetric); err != nil {
		log.Logger.Errorln(err)
	}
	if err := prometheus.Register(reqSz); err != nil {
		log.Logger.Errorln(err)
	}
	if err := prometheus.Register(resSz); err != nil {
		log.Logger.Errorln(err)
	}
	return &Prometheus{
		ReqCnt: reqCntMetric,
		ReqDur: reqDurMetric,
		ReqSz:  reqSz,
		ResSz:  resSz,
	}
}

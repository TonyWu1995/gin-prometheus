package main

import (
	promMiddle "gin-prometheus/cmd/router"
	prom "gin-prometheus/pkg/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	r := gin.Default()
	prometheus := prom.New()
	r.Use(promMiddle.HandlerFunc(prometheus))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"result": "ok",
		})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/ping1", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error",
		})
	})
	r.Run()
}

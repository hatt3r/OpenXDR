package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", Health)
	router.GET("/agents", getAgents)
	router.GET("/events", getEvents)
	router.GET("/alerts", getAlerts)

	// Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}

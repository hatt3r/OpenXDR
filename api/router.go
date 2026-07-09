package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", Health)
	router.GET("/agents", getAgents)
	router.GET("/events", getEvents)
	router.GET("/alerts", getAlerts)
	return router
}

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OPENXDR API is running. ",
	})
}

func getEvents(c *gin.Context) {
	events, err := DB.GetEvents()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, events)
}

func getAlerts(c *gin.Context) {
	alerts, err := DB.GetAlerts()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func getAgents(c *gin.Context) {
	agents, err := DB.GetAgents()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, agents)
}

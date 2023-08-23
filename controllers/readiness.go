package controllers

import (
	"github.com/gin-gonic/gin"
)

func (c *Controller) Readiness(ctx *gin.Context) {
	// Get current server version
	serverVersion, err := c.ClientSet.Discovery().ServerVersion()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "error",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message":     "pong",
		"k8s_version": serverVersion.String(),
	})
	c.Logger.Info().Msg("Readiness probe")
}

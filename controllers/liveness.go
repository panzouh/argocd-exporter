package controllers

import "github.com/gin-gonic/gin"

func (c *Controller) Liveness(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
	c.Logger.Info().Msg("Liveness probe")
}

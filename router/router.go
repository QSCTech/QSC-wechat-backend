package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Load routers and middleware into gin instance
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(mw...)
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Error API not found"})
	})

	return g
}
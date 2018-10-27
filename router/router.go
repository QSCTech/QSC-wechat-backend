package router

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api/sd"
	"git.zjuqsc.com/miniprogram/wechat-backend/router/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Load routers and middleware into gin instance
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Load default middleware
	g.Use(middleware.Options())
	g.Use(middleware.Secure())

	g.Use(mw...)
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Error API not found"})
	})


	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}
	return g
}
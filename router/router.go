package router

import (
	"github.com/gin-gonic/gin"
)

// Load routers and middleware into gin instance
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(mw...)

}
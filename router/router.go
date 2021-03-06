package router

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api/account"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/blackboard"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/ecard"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/exam"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/gpa"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/intlbus"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/login"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/onlineprint"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/schedule"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/sd"
	"git.zjuqsc.com/miniprogram/wechat-backend/api/user"
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

	authGroup := g.Group("/auth")
	{
		authGroup.POST("/login", login.Login)
		authGroup.POST("/bind", login.Bind)
	}

	accountGroup := g.Group("/account")
	accountGroup.Use(middleware.JWTMiddleware())
	{
		accountGroup.GET("", account.GetBindShip)
		accountGroup.POST("/print", account.PrintBind)
		accountGroup.POST("/intl", account.IntlBind)
	}

	userGroup := g.Group("/user")
	userGroup.Use(middleware.JWTMiddleware())
	{
		userGroup.GET("/info", user.Info)
	}

	examGroup := g.Group("/exam")
	examGroup.Use(middleware.JWTMiddleware())
	{
		examGroup.GET("", exam.GetExam)
	}

	gpaGroup := g.Group("/gpa")
	gpaGroup.Use(middleware.JWTMiddleware())
	{
		gpaGroup.GET("", gpa.GetGPA)
		gpaGroup.GET("term", gpa.GetTerm)
	}

	bbGroup := g.Group("/blackboard")
	bbGroup.Use(middleware.JWTMiddleware())
	{
		bbGroup.GET("/alert", blackboard.GetAlert)
	}

	printGroup := g.Group("/print")
	printGroup.Use(middleware.JWTMiddleware())
	{
		printGroup.POST("", onlineprint.Print)
		printGroup.GET("/job", onlineprint.GetJobList)
		printGroup.DELETE("/job", onlineprint.DelJob)
		printGroup.GET("/station", onlineprint.GetStationList)
		printGroup.POST("/bb", onlineprint.PrintBB)
	}

	scheduleGroup := g.Group("/schedule")
	scheduleGroup.Use(middleware.JWTMiddleware())
	{
		scheduleGroup.GET("", schedule.Course)
	}

	intlBusGroup := g.Group("/intlbus")
	intlBusGroup.Use(middleware.JWTMiddleware())
	{
		intlBusGroup.GET("/bus/:date", intlbus.GetBusList)
		intlBusGroup.GET("/book", intlbus.GetBookList)
		intlBusGroup.DELETE("/book", intlbus.DelBook)
		intlBusGroup.GET("/plist", intlbus.GetPlist)
		intlBusGroup.POST("/reserve", intlbus.ReserveBus)
	}

	ecardGroup := g.Group("/ecard")
	ecardGroup.Use(middleware.JWTMiddleware())
	{
		ecardGroup.GET("/balance", ecard.GetBalance)
	}

	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}
	return g
}
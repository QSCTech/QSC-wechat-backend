package main

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/router"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func main() {
	// Settings for Gin
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)


	g := gin.Default()
	router.Load(g)

	g.Run(":8080")
}

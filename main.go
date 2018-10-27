package main

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/config"
	"git.zjuqsc.com/miniprogram/wechat-backend/router"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"io"
	"os"
)

func main() {
	// Init for config and log
	if err := config.Init(""); err != nil {
		panic(err)
	}
	// Settings for Gin
	fgin, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(fgin, os.Stdout)
	gin.SetMode(viper.GetString("RUNMODE"))
	// Create Gin instance
	g := gin.Default()

	// Load router
	router.Load(g)

	log.Infof("Server is listening on %s", viper.GetString("addr"))
	g.Run(viper.GetString("addr"))
}

package main

import (
	"errors"
	"git.zjuqsc.com/miniprogram/wechat-backend/config"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/router"
	"git.zjuqsc.com/miniprogram/wechat-backend/rpc"
	"git.zjuqsc.com/miniprogram/wechat-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"time"
)

func pingServer() error {
	for i := 0; i < 10; i++ {
		resp, err := http.Get( "http://" + viper.GetString("host") + viper.GetString("addr") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}

func healthChecking() {
	if err := pingServer(); err != nil {
		log.Fatal("The router has no response, or it might took too long to start up.", err)
	}
	log.Info("Health checking completes successfully.")
}

func main() {
	// Init for config and log
	if err := config.Init(""); err != nil {
		panic(err)
	}

	// Init for database
	model.DB.Init()
	defer model.DB.Close()

	// Init for gRPC client
	rpc.GRPCClient.Init()
	defer rpc.GRPCClient.Close()

	// Init for Minio client
	service.MinioInit()

	// Settings for Gin
	fgin, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(fgin, os.Stdout)
	gin.SetMode(viper.GetString("RUNMODE"))

	// Create Gin instance
	g := gin.Default()

	// Load router
	router.Load(g)

	// Health checking
	go healthChecking()
	// gRPC checking
	go rpc.PingRPCServer()

	// Start listening
	log.Infof("Server is listening on %s", viper.GetString("addr"))
	g.Run(viper.GetString("addr"))
}

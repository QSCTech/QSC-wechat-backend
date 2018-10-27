package rpc

import (
	"context"
	"git.zjuqsc.com/miniprogram/wechat-backend/protobuf/ZJUPassport"
	"git.zjuqsc.com/miniprogram/wechat-backend/protobuf/ZJUIntl"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type gRPCClient struct {
	Connection *grpc.ClientConn
	Passport ZJUPassport.PassportServiceClient
	Intl	ZJUIntl.IntlServiceClient
}
func (client *gRPCClient) Init()  {
	conn, err := grpc.Dial(viper.GetString("grpc_addr"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err, "did not connect to %s")
	}
	client.Connection = conn
	clientPassport := ZJUPassport.NewPassportServiceClient(conn)
	clientIntl := ZJUIntl.NewIntlServiceClient(conn)
	client.Passport = clientPassport
	client.Intl = clientIntl
	log.Info("gRPC service connect successfully.")
}
func (client *gRPCClient) Close()  {
	client.Connection.Close()
}

var GRPCClient gRPCClient

func PingRPCServer()  {
	resp, err := GRPCClient.Intl.ConnectTest(context.Background(), &ZJUIntl.Infosent{
		Name: "Laphets",
	})
	if err != nil {
		log.Error("gRPC connect error", err)
	}
	log.Infof("ping gRPC success with reply: %s", resp.Message)
}

//func Init() {
//	GRPCClient.Init()
//}
//func Close() {
//	GRPCClient.Close()
//}
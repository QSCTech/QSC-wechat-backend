package service

import (
	"github.com/lexkong/log"
	"github.com/minio/minio-go"
	"github.com/spf13/viper"
)

var MinioClient *minio.Client

func MinioInit() {
	client, err := minio.New(viper.GetString("minio.endpoint"), viper.GetString("minio.accessKeyID"), viper.GetString("minio.secretAccessKey"), viper.GetBool("minio.useSSL"))
	if err != nil {
		log.Fatal("Minio Connect Error", err)
		return
	}
	MinioClient = client
	log.Info("Minio service connect successfully.")
}

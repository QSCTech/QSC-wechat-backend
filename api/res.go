package api

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func Res(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	c.JSON(http.StatusOK, Response{
		Code:code,
		Message:message,
		Data:data,
	})

	//c.AbortWithStatus(200)
}
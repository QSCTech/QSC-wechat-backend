package onlineprint

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/service"
	"github.com/gin-gonic/gin"
)

func GetStationList(c *gin.Context) {
	res, err := service.GetPrintStation()
	if err != nil {
		api.Res(c, errno.ErrPrint, err.Error())
		return
	}

	api.Res(c, nil, res)
}
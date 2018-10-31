package account

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"github.com/gin-gonic/gin"
)

type bindShip struct {
	ZJUid string `json:"ZJUid"`
	INTLid string `json:"INTLid"`
	PRINTid string `json:"PRINTid"`
}

func GetBindShip(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	user, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	res := &bindShip{
		ZJUid: user.ZJUid,
		INTLid: user.INTLid,
		PRINTid: user.PRINTid,
	}

	api.Res(c, nil, res)
}
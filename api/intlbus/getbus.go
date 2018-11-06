package intlbus

import (
	"context"
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/protobuf/INTLUtils"
	"git.zjuqsc.com/miniprogram/wechat-backend/rpc"
	"github.com/gin-gonic/gin"
	"regexp"
)

func GetBusList(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	_, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	date := c.Param("date")
	regExpr := regexp.MustCompile(`^[1-9]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])$`)
	if regRes := regExpr.MatchString(date); regRes == false {
		api.Res(c, errno.ErrDateFormat, regRes)
		return
	}

	rpcReq := &INTLUtils.BusReq{
		ZJUid: ZJUid,
		Date: date,
	}
	resp, err := rpc.GRPCClient.IntlBus.GetBusList(context.Background(), rpcReq)
	if resp.Status.Code != 200 {
		api.Res(c, errno.ErrBusGet, resp.Status.Info)
		return
	}

	api.Res(c, nil, resp.Buslist)
}
package intlbus

import (
	"context"
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/protobuf/INTLUtils"
	"git.zjuqsc.com/miniprogram/wechat-backend/rpc"
	"github.com/gin-gonic/gin"
)

func GetPlist(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	_, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	rpcReq := &INTLUtils.PlistReq{
		ZJUid: ZJUid,
	}
	resp, err := rpc.GRPCClient.IntlBus.GetPlist(context.Background(), rpcReq)
	if resp.Status.Code != 200 {
		api.Res(c, errno.ErrBusGet, resp.Status.Info)
		return
	}

	api.Res(c, nil, resp.Plist)
}
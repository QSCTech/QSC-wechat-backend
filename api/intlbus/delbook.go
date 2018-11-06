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

type delbookRequest struct {
	Bookid int32 `json:"bookid"`
}
func DelBook(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	_, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	req := delbookRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Res(c, errno.ErrBind, err.Error())
		return
	}

	rpcReq := &INTLUtils.DelBookReq{
		ZJUid: ZJUid,
		Bookid: req.Bookid,
	}
	resp, err := rpc.GRPCClient.IntlBus.DelBook(context.Background(), rpcReq)

	if resp.Status.Code != 200 {
		api.Res(c, errno.ErrBusPatch, resp.Status.Info)
		return
	}

	api.Res(c, nil, nil)
}
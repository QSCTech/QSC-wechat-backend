package blackboard

import (
	"context"
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/crypt"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/protobuf/ZJUIntl"
	"git.zjuqsc.com/miniprogram/wechat-backend/rpc"
	"github.com/gin-gonic/gin"
)

func GetAlert(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	user, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	decryptedPass, err := crypt.Decrypt(user.Password)
	if err != nil {
		api.Res(c, errno.ErrDecrypt, err.Error())
		return
	}

	rpcReq := &ZJUIntl.AlertListReq{
		Username: ZJUid,
		Password: decryptedPass,
	}
	resp, err := rpc.GRPCClient.BlackBoard.GetAlertList(context.Background(), rpcReq)
	if err != nil {
		api.Res(c, errno.ErrBBAlert, err.Error())
		return
	}
	if resp.Status.Code != 200 {
		api.Res(c, errno.ErrBBAlert, resp.Status.Info)
		return
	}

	api.Res(c, nil, resp.Alertlist)
}
package gpa

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

func GetGPA(c *gin.Context) {
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

	term := c.Query("term")

	rpcReq := &ZJUIntl.GPAInfoReq{
		Username: ZJUid,
		Password: decryptedPass,
		Term: term,
	}
	resp, err := rpc.GRPCClient.Intl.GetGPAInfo(context.Background(), rpcReq)
	if err != nil {
		api.Res(c, errno.ErrBBAlert, err.Error())
		return
	}
	if resp.Status.Code != 200 {
		api.Res(c, errno.ErrBBAlert, resp.Status.Info)
		return
	}

	api.Res(c, nil, resp)
}
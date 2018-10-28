package schedule

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

func Course(c *gin.Context) {
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

	// Get schedule from gRPC service
	res, err := rpc.GRPCClient.Intl.GetCourse(context.Background(), &ZJUIntl.User{
		Username: ZJUid,
		Password: decryptedPass,
	})
	if err != nil {
		api.Res(c, errno.ErrRPC, err.Error())
		return
	}
	if res.Status != "SUCCESS" {
		api.Res(c, errno.ErrRPC, err.Error())
		return
	}

	api.Res(c, nil, res.Course)
}
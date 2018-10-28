package user

import (
	"context"
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/protobuf/ZJUPassport"
	"git.zjuqsc.com/miniprogram/wechat-backend/rpc"
	"github.com/gin-gonic/gin"
)

//type infoResponse struct {
//
//}
func Info(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	_, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}
	res, err := rpc.GRPCClient.Passport.GetAllInfo(context.Background(), &ZJUPassport.BasicRequest{
		ZJUid:ZJUid,
		Password: "",
		Mode: ZJUPassport.Mode_CACHE,
	})
	if err != nil {
		api.Res(c, errno.ErrGetStudent, err.Error())
		return
	}
	if res.Status.Code != 200 {
		api.Res(c, errno.ErrGetStudent, res.Status.Info)
		return
	}
	api.Res(c, nil, res.Student)
}
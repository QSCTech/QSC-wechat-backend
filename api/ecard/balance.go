package ecard

import (
	"context"
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/crypt"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/protobuf/ZJUPassport"
	"git.zjuqsc.com/miniprogram/wechat-backend/rpc"
	"github.com/gin-gonic/gin"
)

type cardResponse struct {
	CardID string `json:"card_id"`
	Balance float32 `json:"balance"`
	Daycostamt float32 `json:"daycostamt"`
}
func GetBalance(c *gin.Context) {
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
	res, err := rpc.GRPCClient.Passport.GetAllInfo(context.Background(), &ZJUPassport.BasicRequest{
		ZJUid: ZJUid,
		Password: decryptedPass,
		Mode: ZJUPassport.Mode_INSTANT,
	})

	if err != nil {
		api.Res(c, errno.ErrRPC, err.Error())
		return
	}
	if res.Status.Code != 200 {
		api.Res(c, errno.ErrRPC, res.Status.Info)
		return
	}

	response := &cardResponse{
		CardID: res.Student.CardID,
		Balance: res.Student.Balance,
		Daycostamt: res.Student.Daycostamt,
	}
	api.Res(c, nil, response)
}
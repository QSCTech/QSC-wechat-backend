package login

import (
	"context"
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/crypt"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/token"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/wechat"
	"git.zjuqsc.com/miniprogram/wechat-backend/protobuf/ZJUPassport"
	"git.zjuqsc.com/miniprogram/wechat-backend/rpc"
	"github.com/gin-gonic/gin"
)

type bindRequest struct {
	ZJUid string `json:"ZJUid"`
	Password string `json:"password"`
	Code string `json:"code"`
}
type bindResponse struct {
	AccessToken string `json:"access_token"`
}
func Bind(c *gin.Context)  {
	req := bindRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Res(c, errno.ErrBind, err.Error())
		return
	}

	// First checking ZJUid and password (get personal info by passport)
	res, err := rpc.GRPCClient.Passport.GetAllInfo(context.Background(), &ZJUPassport.BasicRequest{
		ZJUid: req.ZJUid,
		Password: req.Password,
		Mode: ZJUPassport.Mode_INSTANT,
	})
	if err != nil {
		api.Res(c, errno.ErrGetStudent, err.Error())
		return
	}
	if res.Status.Code != 200 {
		api.Res(c, errno.ErrGetStudent, res.Status.Info)
		return
	}

	// Student checking successfully, bind with wechatOpenID
	session, err := wechat.Code2Session(req.Code)
	if err != nil {
		api.Res(c, errno.ErrWechat, err.Error())
		return
	}

	// Delete existing bindingship
	if err := model.DeleteZJU(req.ZJUid); err != nil {
		api.Res(c, errno.ErrDeBind, err.Error())
		return
	}
	if err := model.DeleteWechat(session.OpenID); err != nil {
		api.Res(c, errno.ErrDeBind, err.Error())
		return
	}

	encryptedPass, err := crypt.Encrypt(req.Password)
	if err != nil {
		api.Res(c, errno.ErrEncrypt, err.Error())
		return
	}
	// Insert new binding ship
	user := &model.UserModel{
		ZJUid: req.ZJUid,
		Password: encryptedPass,
		WechatOpenID: session.OpenID,
		WechatSessionID: session.SessionKey,
	}
	if err := user.Create(); err != nil {
		api.Res(c, errno.DBError, err.Error())
		return
	}

	// Return token
	JWT, err := token.Sign(token.Context{
		ZJUid: req.ZJUid,
	}, "")
	if err != nil {
		api.Res(c, errno.ErrToken, err.Error())
		return
	}

	api.Res(c, nil, &bindResponse{
		AccessToken: JWT,
	})
}
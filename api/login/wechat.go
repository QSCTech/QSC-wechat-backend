package login

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/token"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/wechat"
	"github.com/gin-gonic/gin"
)

// This function is used to directly using OpenID to login
// Which is different from bind
type loginRequest struct {
	Code string `json:"code"  binding:"required"`
}
func Login(c *gin.Context)  {
	req := loginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Res(c, errno.ErrBind, err.Error())
		return
	}

	session, err := wechat.Code2Session(req.Code)
	if err != nil {
		api.Res(c, errno.ErrWechat, err.Error())
		return
	}

	user, err := model.GetUserByWechatID(session.OpenID)
	if err != nil {
		api.Res(c, errno.ErrNoBindingUser, err.Error())
		return
	}

	log4Save := model.LogModel{
		ZJUid: user.ZJUid,
		IP: c.ClientIP(),
		URL: c.Request.RequestURI,
		UA: c.GetHeader("User-Agent"),
	}
	go log4Save.Create()

	// Return token
	JWT, err := token.Sign(token.Context{
		ZJUid: user.ZJUid,
	}, "")
	if err != nil {
		api.Res(c, errno.ErrToken, err.Error())
		return
	}

	api.Res(c, nil, &bindResponse{
		AccessToken: JWT,
	})
}
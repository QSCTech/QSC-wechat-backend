package login

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/wechat"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
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

	session := wechat.Code2Session(req.Code)

	log.Infof("%+v", session)

}
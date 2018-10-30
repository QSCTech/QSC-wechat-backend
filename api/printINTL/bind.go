package printINTL

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/crypt"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/service"
	"github.com/gin-gonic/gin"
)

type bindRequest struct {
	Password string `json:"password"`
}
func Bind(c *gin.Context) {
	req := bindRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Res(c, errno.ErrBind, err.Error())
		return
	}

	ZJUid := c.GetString("ZJUid")
	user, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	// Validate for password
	if err := service.CheckPrintPassword(ZJUid, req.Password); err != nil {
		api.Res(c, errno.ErrPassword, err.Error())
		return
	}

	passEncrypted, err := crypt.Encrypt(req.Password)
	if err != nil {
		api.Res(c, errno.ErrEncrypt, err.Error())
		return
	}
	user.PasswordPRINT = passEncrypted
	if err := user.Save(); err != nil {
		api.Res(c, errno.ErrBind, err.Error())
		return
	}

	api.Res(c, nil, nil)
}
package onlineprint

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/crypt"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/service"
	"github.com/gin-gonic/gin"
)

func Print(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	user, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		api.Res(c, errno.ErrFileRead, "请上传有效的文件.")
		return
	}

	passwordDecrepted, err := crypt.Decrypt(user.PasswordPRINT)
	if err != nil {
		api.Res(c, errno.ErrDecrypt, err.Error())
		return
	}

	paperid := c.DefaultPostForm("paperid", "A4")
	color := c.DefaultPostForm("color", "0")
	double := c.DefaultPostForm("double", "dupnone")
	copies := c.DefaultPostForm("copies", "1")

	fields := map[string]string{
		"paperid": paperid,
		"color": color,
		"double": double,
		"copies": copies,
	}

	if err := service.Print(ZJUid, passwordDecrepted, fields, file); err != nil {
		api.Res(c, errno.ErrPrint, err.Error())
		return
	}

	api.Res(c, nil, nil)
}
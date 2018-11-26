package onlineprint

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/crypt"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/service"
	"github.com/gin-gonic/gin"
)

func GetJobList(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	user, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	passwordDecrepted, err := crypt.Decrypt(user.PasswordPRINT)
	if err != nil {
		api.Res(c, errno.ErrDecrypt, err.Error())
		return
	}

	res, err := service.GetPrintJob(ZJUid, passwordDecrepted)
	if err != nil {
		api.Res(c, errno.ErrPrint, err.Error())
		return
	}

	api.Res(c, nil, res)
}

type deljobRequest struct {
	JobId int32 `json:"job_id" binding:"required"`
}
func DelJob(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	user, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	req := deljobRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Res(c, errno.ErrBind, err.Error())
		return
	}

	passwordDecrepted, err := crypt.Decrypt(user.PasswordPRINT)
	if err != nil {
		api.Res(c, errno.ErrDecrypt, err.Error())
		return
	}

	if err := service.DeletePrintJob(ZJUid, passwordDecrepted, req.JobId); err != nil {
		api.Res(c, errno.ErrPrint, err.Error())
		return
	}

	api.Res(c,nil, nil)
}
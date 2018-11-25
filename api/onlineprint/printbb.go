package onlineprint

import (
	"context"
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/model"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/crypt"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/protobuf/ZJUIntl"
	"git.zjuqsc.com/miniprogram/wechat-backend/rpc"
	"git.zjuqsc.com/miniprogram/wechat-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
)

type printBBRequest struct {
	FileUrl string `json:"file_url" binding:"required"`
	PaperId string `json:"paper_id"`
	Color string `json:"color"`
	Double string `json:"double"`
	Copies int32 `json:"copies"`
}
func PrintBB(c *gin.Context) {
	ZJUid := c.GetString("ZJUid")
	user, err := model.GetUserByZJUid(ZJUid)
	if err != nil {
		api.Res(c, errno.ErrUserNotFound, err.Error())
		return
	}

	req := &printBBRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		api.Res(c, errno.ErrBind, err.Error())
		return
	}

	decryptedPass, err := crypt.Decrypt(user.Password)
	if err != nil {
		api.Res(c, errno.ErrDecrypt, err.Error())
		return
	}

	rpcReq := &ZJUIntl.FileEtagReq{
		FileUrl: req.FileUrl,
		Username: ZJUid,
		Password: decryptedPass,
	}

	resp, err := rpc.GRPCClient.BlackBoard.GetFileEtag(context.Background(), rpcReq)
	if err != nil {
		api.Res(c, errno.ErrBBFile, err.Error())
		return
	}
	if resp.Status.Code != 200 {
		api.Res(c, errno.ErrBBFile, resp.Status.Info)
		return
	}

	// Then fetch file and upload


	object, err := service.MinioClient.GetObject("bbfiles", resp.FilePath, minio.GetObjectOptions{})
	if err != nil {
		api.Res(c, errno.ErrOSS, err.Error())
		return
	}

	objInfo, err := object.Stat()
	if err != nil {
		api.Res(c, errno.ErrOSS, err.Error())
		return
	}

	fields := map[string]string{
		"paperid": "A4",
		"color": "0",
		"double": "dupnone",
		"copies": "1",
	}

	printPasswordDecrepted, err := crypt.Decrypt(user.PasswordPRINT)
	if err != nil {
		api.Res(c, errno.ErrDecrypt, err.Error())
		return
	}

	if err := service.PrintFromMinio(ZJUid, printPasswordDecrepted, fields, object, objInfo.Key); err != nil {
		api.Res(c, errno.ErrPrint, err.Error())
		return
	}

	api.Res(c, nil, resp.FilePath)
}
package middleware

import (
	"git.zjuqsc.com/miniprogram/wechat-backend/api"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/errno"
	"git.zjuqsc.com/miniprogram/wechat-backend/pkg/token"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		JWTpayload, err := token.ParseRequest(c)
		if err != nil {
			api.Res(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		c.Set("ZJUid", JWTpayload.ZJUid)
		c.Set("INTLid", JWTpayload.INTLid)
		c.Next()
	}
}
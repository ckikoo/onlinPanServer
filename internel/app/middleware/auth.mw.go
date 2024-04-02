package middleware

import (
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/pkg/auth"
	"onlineCLoud/pkg/contextx"
	"strings"

	"github.com/gin-gonic/gin"
)

// TODO: ADD support admin
func wrapUserAuthContext(c *gin.Context, userID string, email string, admin ...bool) {
	ctx := contextx.NewUserEmail(c.Request.Context(), email)
	ctx = contextx.NewUserID(ctx, userID)
	ctx = contextx.NewUUID(ctx)
	c.Request = c.Request.WithContext(ctx)
}

// TODO ADD support for admin diff
// 权限中间件 并提取token存放的关键数据到上下文
func AuthMiddleware(a auth.Auther, skipper ...SkipperFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if SkipHandler(ctx, skipper...) {
			ctx.Next()
			return
		}

		tokenUserEmain, err := a.ParseUserEmail(ctx.Request.Context(), ginx.GetToken(ctx))
		if err != nil {
			if err == auth.ErrInvalidToken {
				ginx.ResFailWithMessage(ctx, "请重新登录")
				ctx.Abort()
				return
			}
			ginx.ResFailWithMessage(ctx, err.Error())
			ctx.Abort()
			return
		}

		idx := strings.Index(tokenUserEmain, " ")
		if idx == -1 {
			ginx.ResFailWithMessage(ctx, "用户已过期请重新登录")
			ctx.Abort()
			return
		}
		userID := tokenUserEmain[:idx]
		wrapUserAuthContext(ctx, userID, tokenUserEmain[idx+1:])
		ctx.Next()
	}
}

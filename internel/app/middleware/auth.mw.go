package middleware

import (
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/pkg/auth"
	"onlineCLoud/pkg/contextx"

	"github.com/gin-gonic/gin"
)

// 校验用户新的中间层
func AuthMiddleware(a auth.Auther, skipper SkipperFunc, admin SkipperFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if SkipHandler(ctx, skipper) {
			ctx.Next()
			return
		}

		if contextx.FromMiddle(ctx.Request.Context()) != "" {
			ginx.ResFailWithMessage(ctx, contextx.FromMiddle(ctx.Request.Context()))
			ctx.Abort()

			return
		}

		if SkipHandler(ctx, admin) {
			ctx.Next()
			return
		}

		ctx.Next()
	}
}

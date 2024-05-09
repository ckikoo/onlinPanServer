package middleware

import (
	"fmt"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/pkg/auth"
	"onlineCLoud/pkg/contextx"

	"github.com/gin-gonic/gin"
)

// 校验用户新的中间层
func AuthMiddleware(a auth.Auther, skipper ...SkipperFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if SkipHandler(ctx, skipper...) {
			ctx.Next()
			return
		}

		if contextx.FromMiddle(ctx.Request.Context()) != "" {
			ginx.ResNeedReload(ctx)
			fmt.Println("1111111111")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

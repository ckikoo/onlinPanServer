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

func wrapMIddleReason(c *gin.Context, res string) {
	ctx := contextx.NewMiddle(c.Request.Context(), res)

	c.Request = c.Request.WithContext(ctx)
}

func UserInfo(a auth.Auther) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		tokenUserEmain, err := a.ParseUserEmail(ctx.Request.Context(), ginx.GetToken(ctx))
		if err != nil {

			if err == auth.ErrInvalidToken {
				wrapMIddleReason(ctx, "请重新登录")
				ctx.Next()
				return
			}
			wrapMIddleReason(ctx, err.Error())
			ctx.Next()
			return
		}

		idx := strings.Index(tokenUserEmain, " ")
		if idx == -1 {
			wrapMIddleReason(ctx, "请重新登录")
			ctx.Next()
			return
		}

		userID := tokenUserEmain[:idx]
		wrapUserAuthContext(ctx, userID, tokenUserEmain[idx+1:])
		ctx.Next()
	}
}

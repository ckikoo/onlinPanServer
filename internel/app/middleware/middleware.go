package middleware

import (
	"fmt"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/pkg/contextx"
	"onlineCLoud/pkg/errors"

	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NoMethodHander() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ginx.ResJson(ctx, 404, map[string]interface{}{}, errors.ErrBadRequest, "false")
	}
}
func NoRouteHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ginx.ResJson(ctx, 404, map[string]interface{}{}, errors.ErrBadRequest, "false")
	}
}

type SkipperFunc func(*gin.Context) bool

func AllowAdminSkipper(prefix ...string) SkipperFunc {
	return func(ctx *gin.Context) bool {

		path := ctx.Request.URL.Path

		pathLen := len(path)
		adminPrex := "/api/admin"
		adminLen := len(adminPrex)

		//  排除前缀不属于的
		if pathLen < adminLen || path[:adminLen] != adminPrex {
			return true
		}

		session := sessions.Default(ctx)
		isAdmin := session.Get("pri").(string)

		for _, p := range prefix {
			if p1 := len(p); pathLen >= p1 {
				admin := contextx.GetAdmin(ctx.Request.Context())
				if admin == "1" && isAdmin == "admin" {
					return true
				}
			}
		}
		return false
	}
}

func AllowPathPrefixSkipper(prefix ...string) SkipperFunc {
	return func(ctx *gin.Context) bool {
		path := ctx.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefix {
			if p1 := len(p); pathLen >= p1 && path[:p1] == p {
				return true
			}
		}
		return false
	}
}

func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

func AllowMethodAndPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := JoinRouter(c.Request.Method, c.Request.URL.Path)
		pathLen := len(path)
		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}
func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}

func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

func EmptyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func PrintUrlRequest() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		begin := time.Now().UTC()
		fmt.Println(ctx.Request.URL.Path)
		ctx.Next()
		diff := time.Since(begin).Microseconds()
		fmt.Println(ctx.Request.URL.Path + " 请求结束，处理时间：" + strconv.FormatInt(diff, 10) + " 毫秒")

	}
}

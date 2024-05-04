package middleware

import (
	"github.com/gin-gonic/gin"
)

func PrintUrlRequest() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// begin := time.Now().UTC()
		ctx.Next()
		// diff := time.Since(begin).Milliseconds()

	}
}

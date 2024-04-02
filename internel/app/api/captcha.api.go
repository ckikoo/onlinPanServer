package api

import (
	"fmt"
	"onlineCLoud/internel/app/config"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"

	"github.com/gin-gonic/gin"
)

func GenerateCaptcha(ctx *gin.Context) {
	Type := ctx.Query("type")
	fmt.Printf("Type: %v\n", Type)
	cfg := config.C.Captcha
	service.Captcha(ctx, Type, cfg.Length, cfg.Width, cfg.Height)
}

func Verfify(ctx *gin.Context) {
	value := ctx.PostForm("code")
	Type := ctx.PostForm("type")

	if service.CaptchaVerify(ctx, Type, value) {
		ginx.ResOkWithMessage(ctx, "验证码正确")
	} else {
		ginx.ResFailWithMessage(ctx, "验证码错误")
	}
}

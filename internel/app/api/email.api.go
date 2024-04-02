package api

import (
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/util/random"

	"github.com/gin-gonic/gin"
)

func SendEmail(c *gin.Context) {
	Type := c.PostForm("type")
	checkCode := c.PostForm("checkCode")
	email := c.PostForm("email")

	if Type == "1" && service.CaptchaVerify(c, Type, checkCode) {
		service.Email(c.Request.Context(), email, random.GetRandom(4))
		ginx.ResOk(c)
		return
	}
	ginx.ResFailWithMessage(c, "验证码错误")
}

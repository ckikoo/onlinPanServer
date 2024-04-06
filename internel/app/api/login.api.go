package api

import (
	"encoding/json"
	"fmt"

	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"

	"github.com/gin-gonic/gin"
)

type LoginAPI struct {
	LoginSrv *service.LoginSrv
}

func (a *LoginAPI) formatTokenUserID(userID string, userName string) string {
	return fmt.Sprintf("%s %s", userID, userName)
}

func (a *LoginAPI) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item user.User
	if err := ginx.ParseForm(c, &item); err != nil {
		ginx.ResFailWithMessage(c, "数据错误")
		return
	}

	checkCode := c.PostForm("checkCode")
	if !service.CaptchaVerify(c, "0", checkCode) {
		ginx.ResFailWithMessage(c, "验证码错误")
		return
	}
	user_id, err := a.LoginSrv.Login(ctx, item.Email, item.Password)

	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	tokenMap, _ := a.LoginSrv.GenerateToken(ctx, a.formatTokenUserID(user_id, item.Email))
	v, _ := json.Marshal(tokenMap)

	ginx.ResOkWithData(c, string(v))
}

func (a *LoginAPI) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var item user.User
	if err := ginx.ParseForm(c, &item); err != nil {
		ginx.ResFailWithMessage(c, "数据错误")
		return
	}

	ok, err := a.LoginSrv.Register(ctx, item)

	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	} else if !ok {
		ginx.ResFailWithMessage(c, "注册失败")
		return
	} else {
		ginx.ResOkWithMessage(c, "注册成功")
	}

}

func (a *LoginAPI) ResetPasswd(c *gin.Context) {
	ctx := c.Request.Context()
	email := c.PostForm("email")
	password := c.PostForm("password")
	checkCode := c.PostForm("checkCode")
	emailCode := c.PostForm("emailCode")
	if !service.CaptchaVerify(c, "0", checkCode) {
		ginx.ResFailWithMessage(c, "验证码错误")
		return
	}
	if ok, _ := service.CheckEmail(ctx, email, emailCode); !ok {
		ginx.ResFailWithMessage(c, "邮箱验证码错误")
		return
	}

	go service.DeleteEmail(ctx, email)

	var item *user.User
	item = a.LoginSrv.FindOneByName(ctx, email)
	if item.Password == "" {
		ginx.ResFailWithMessage(c, "用户不存在")
		return
	}
	ok, err := a.LoginSrv.ResetPasswd(ctx, email, password)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	if ok {
		ginx.ResOkWithMessage(c, "密码重置成功")
	} else {
		ginx.ResFailWithMessage(c, "密码重置失败")
	}

}

func (a *LoginAPI) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	userEmail := contextx.FromUserEmail(ctx)
	if userEmail != "" {
		_ = a.LoginSrv.DestoryToken(ctx, ginx.GetToken(c))
	}
	c.SetCookie("userinfo", "", -1, "/", "localhost", false, true)
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)

	ginx.ResOk(c)
}

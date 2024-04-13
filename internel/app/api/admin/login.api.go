package admin

import (
	"encoding/json"
	"fmt"

	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminLoginAPI struct {
	LoginSrv *service.LoginSrv
}

func (a *AdminLoginAPI) formatTokenUserID(userID string, userName string) string {
	return fmt.Sprintf("%s %s %d", userID, userName, 1)
}

func (a *AdminLoginAPI) Login(c *gin.Context) {
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
	user_id, err := a.LoginSrv.Login(ctx, item.Email, item.Password, true)

	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	session := sessions.Default(c)
	session.Set("pri", "admin")
	session.Save()

	tokenMap, _ := a.LoginSrv.GenerateToken(ctx, a.formatTokenUserID(user_id, item.Email))
	v, _ := json.Marshal(tokenMap)

	ginx.ResOkWithData(c, string(v))
}

func (a *AdminLoginAPI) ResetPasswd(c *gin.Context) {
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

func (a *AdminLoginAPI) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	userEmail := contextx.FromUserEmail(ctx)
	if userEmail != "" {
		_ = a.LoginSrv.DestoryToken(ctx, ginx.GetToken(c))
	}

	c.SetCookie("userinfo", "", -1, "/", "localhost", false, true)
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	session := sessions.Default(c)
	session.Delete("pri")
	session.Save()

	ginx.ResOk(c)
}

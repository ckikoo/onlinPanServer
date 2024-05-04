package api

import (
	"fmt"
	"io"
	"log"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	UserSrv *service.UserSrv
}

func (a *UserAPI) GetInfo(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserSrv.GetInfo(ctx, contextx.FromUserEmail(ctx))
	if err != nil {
		ginx.ResFailWithMessage(c, "获取用户信息失败")
		return
	}
	ginx.ResOkWithData(c, item)
}

func (a *UserAPI) UpdateUserAvatar(c *gin.Context) {
	const (
		MaxAvatarSize = 10485760
	)
	ctx := c.Request.Context()
	fh, err := c.FormFile("avatar")
	if err != nil {
		ginx.ResFail(c)
		return
	}

	fileExt := path.Ext(fh.Filename)
	if fileExt != ".jpg" && fileExt != ".png" {
		ginx.ResFailWithMessage(c, "仅支持 .jpg 和 .png 格式的文件")
		return
	}

	if fh.Size > MaxAvatarSize {
		ginx.ResFailWithMessage(c, "文件大小超过限制")
		return
	}

	filename := fmt.Sprintf("img/%v", contextx.FromUserID(ctx)+path.Ext(fh.Filename))

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Default().Println(err)
		ginx.ResFail(c)
		return
	}
	defer f.Close()
	src, err := fh.Open()
	if err != nil {
		log.Default().Println(err)
		ginx.ResFail(c)
		return
	}
	defer src.Close()

	_, err = io.Copy(f, src)
	if err != nil {
		log.Default().Println(err)
		ginx.ResFailWithMessage(c, "文件保存失败")
		return
	}

	err = a.UserSrv.UpdateUserAvatar(ctx, contextx.FromUserEmail(ctx), filename)
	if err != nil {
		log.Default().Println(err)

		ginx.ResFail(c)
		return
	}

	ginx.ResOk(c)
}
func (a *UserAPI) GetUserSpace(c *gin.Context) {
	ctx := c.Request.Context()

	item := a.UserSrv.GetUserSpace(ctx, contextx.FromUserEmail(ctx))
	ginx.ResOkWithData(c, item)

}

func (a *UserAPI) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	oldPassword := c.PostForm("oldPassword")
	newPassword := c.PostForm("password")
	err := a.UserSrv.UpdatePassword(ctx, contextx.FromUserEmail(ctx), oldPassword, newPassword)
	if err != nil {
		log.Println(err)
		ginx.ResFail(c)
		return
	}

	ginx.ResOk(c)

}
func (a *UserAPI) GetUserAvatar(c *gin.Context) {

	uid := c.Param("user")
	if uid == "" {
		ginx.ResFail(c)
	}
	err := a.UserSrv.GetUserAvatar(c.Writer, c.Request, uid)
	if err != nil {
		ginx.ResFailWithMessage(c, "获取用户信息失败")
		return
	}
}

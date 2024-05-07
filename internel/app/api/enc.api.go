package api

import (
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type EncAPI struct {
	EncSrv  *service.EncSrv
	FileSrv *service.FileSrv
}

func (api *EncAPI) checkSession(c *gin.Context) int {
	sesion := sessions.Default(c)
	if check := sesion.Get("pass"); check == nil {
		return 0
	}

	return 1

}

func (api *EncAPI) AddFile(c *gin.Context) {
	ctx := c.Request.Context()

	fileIds := c.PostForm("fileIds")
	if fileIds == "" {
		ginx.ResFail(c)
		return
	}

	err := api.FileSrv.UpdateFileSecure(ctx, contextx.FromUserID(ctx), fileIds, true)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOk(c)
}
func (api *EncAPI) InitPassword(c *gin.Context) {
	ctx := c.Request.Context()
	pass := c.PostForm("encpass")
	if pass == "" {
		ginx.ResFail(c)
		return
	}

	err := api.EncSrv.InitPassWord(ctx, contextx.FromUserEmail(ctx), pass)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	session := sessions.Default(c)
	session.Set("pass", 1)
	session.Save()

	ginx.ResOk(c)
}

func (api *EncAPI) CheckPassword(c *gin.Context) {
	ctx := c.Request.Context()
	pass := c.PostForm("encpass")
	if pass == "" {
		ginx.ResFail(c)
		return
	}

	f, err := api.EncSrv.CheckPassword(ctx, contextx.FromUserEmail(ctx), pass)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	if f {
		session := sessions.Default(c)
		session.Set("pass", 1)
		session.Save()
		ginx.ResOk(c)
	} else {
		ginx.ResFailWithMessage(c, "密码错误")
	}
}

func (api *EncAPI) CheckEnc(c *gin.Context) {
	ctx := c.Request.Context()

	if api.checkSession(c) == 1 {
		ginx.ResOkWithData(c, 2)
		return
	}

	f, err := api.EncSrv.CheckPassword(ctx, contextx.FromUserEmail(ctx), "")
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	if f {
		ginx.ResOkWithData(c, 0)
	} else {
		ginx.ResOkWithData(c, 1)
	}
}
func (api *EncAPI) LoadencList(c *gin.Context) {
	if api.checkSession(c) == 0 {
		ginx.ResNeedReload(c)
		return
	}

	ctx := c.Request.Context()
	pageNo := c.DefaultPostForm("pageNo", "1")
	pageSize := c.DefaultPostForm("pageSize", "10")
	fileNameFuzzy := c.DefaultPostForm("fileNameFuzzy", "")
	if pageNo == "" {
		pageNo = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	PageNo, err := strconv.ParseInt(pageNo, 10, 64)
	if err != nil {
		ginx.ResFailWithMessage(c, "获取数据失败")
		return
	}
	PageSize, err := strconv.ParseInt(pageSize, 10, 64)

	query := schema.RequestFileListPage{
		PageParams: schema.PageParams{
			PageNo:   int(PageNo),
			PageSize: int(PageSize),
		},
		FileNameFuzzy: fileNameFuzzy,
		Secure:        true,
	}

	if err != nil {
		ginx.ResFailWithMessage(c, "获取数据失败")
		return
	}
	res, err := api.FileSrv.LoadListFiles(c, contextx.FromUserID(ctx), &query)
	if err != nil {
		ginx.ResFailWithMessage(c, "获取数据失败")
		return
	}

	ginx.ResOkWithData(c, res)
}
func (api *EncAPI) DelFile(c *gin.Context) {
	if api.checkSession(c) == 0 {
		ginx.ResNeedReload(c)
		return
	}
}
func (api *EncAPI) RecoverFile(c *gin.Context) {
	ctx := c.Request.Context()
	if api.checkSession(c) == 0 {
		ginx.ResNeedReload(c)
		return
	}
	fileIds := c.PostForm("fileIds")

	err := api.FileSrv.UpdateFileSecure(ctx, contextx.FromUserID(ctx), fileIds, false)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
	} else {
		ginx.ResOk(c)
	}
}

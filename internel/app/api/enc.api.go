package api

import (
	"fmt"
	"onlineCLoud/internel/app/config"
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

	err = api.EncSrv.EnFilePass(ctx, contextx.FromUserID(ctx), fileIds)
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
	if api.checkSession(c) == 0 && config.C.RunMode != "DEBUG" {
		ginx.ResNeedReload(c)
		return
	}

	ctx := c.Request.Context()
	pageNo := c.DefaultPostForm("pageNo", "1")
	pageSize := c.DefaultPostForm("pageSize", "10")
	fileNameFuzzy := c.DefaultPostForm("fileNameFuzzy", "")
	filepid := c.PostForm("filePid")
	if filepid == "" {
		ginx.ResFail(c)
		return
	}
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
		FilePid: filepid,
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
	ctx := c.Request.Context()
	fileids := c.PostForm("fileIds")
	err := api.FileSrv.DelFiles(c, contextx.FromUserID(ctx), fileids, true)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		ginx.ResFailWithMessage(c, "删除失败")
		return
	}

	ginx.ResOk(c)
}

func (api *EncAPI) RecoverFile(c *gin.Context) {
	ctx := c.Request.Context()
	if api.checkSession(c) == 0 {
		ginx.ResNeedReload(c)
		return
	}

	fileIds := c.PostForm("fileIds")

	err := api.FileSrv.UpdateFileSecure(ctx, contextx.FromUserID(ctx), fileIds, false) // 修改文件状态
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
	} else {
		ginx.ResOk(c)
	}

}

func (api *EncAPI) DelFiles(c *gin.Context) {
	ctx := c.Request.Context()

	fileIds := c.PostForm("fileIds")
	err := api.FileSrv.DeleteFiles(ctx, contextx.FromUserID(ctx), fileIds)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOk(c)
}

func (api *EncAPI) NewFoloder(c *gin.Context) {
	ctx := c.Request.Context()

	filePid := c.PostForm("filePid")
	fileName := c.PostForm("fileName")

	info, err := api.FileSrv.NewFoloder(c, contextx.FromUserID(ctx), filePid, fileName, true)
	if err != nil {
		ginx.ResFailWithMessage(c, "创建失败")
		return
	}
	ginx.ResOkWithData(c, info)
}
func (api *EncAPI) LoadAllFolder(c *gin.Context) {
	ctx := c.Request.Context()
	filePid := c.PostForm("filePid")
	currentFileIds := c.PostForm("currentFileIds")
	if filePid == "" {
		ginx.ResFail(c)
	}
	files, err := api.FileSrv.LoadAllFolder(ctx, contextx.FromUserID(ctx), filePid, currentFileIds, true)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, files)
}

func (api *EncAPI) ChangeFileFolder(c *gin.Context) {

	ctx := c.Request.Context()

	fileIds := c.Request.FormValue("fileIds")
	filePid := c.Request.FormValue("filePid")
	if filePid == "" || fileIds == "" {
		ginx.ResFail(c)
	}
	err := api.FileSrv.ChangeFileFolder(ctx, contextx.FromUserID(ctx), fileIds, filePid, true)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	ginx.ResOk(c)
}

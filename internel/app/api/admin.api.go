package api

import (
	"fmt"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminApi struct {
	AdminSrv *service.AdminSrv
}

func (a *AdminApi) LoadUserList(c *gin.Context) {
	ctx := c.Request.Context()
	pageNo := c.Request.PostFormValue("pageNo")
	pageSize := c.Request.PostFormValue("pageSize")
	nickNameFuzzy := c.PostForm("nickNameFuzzy")
	status := c.DefaultPostForm("status", "*")
	if pageNo == "" && pageSize == "" {
		pageNo = "1"
		pageSize = "20"
	}

	PageNo, err := strconv.ParseInt(pageNo, 10, 64)
	if err != nil {
		ginx.ResFail(c)
	}
	PageSize, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		ginx.ResFail(c)
	}

	res, err := a.AdminSrv.LoadUserList(ctx, int(PageNo), int(PageSize), nickNameFuzzy, status)
	if err != nil {
		ginx.ResFail(c)
	}

	ginx.ResOkWithData(c, res)
}

func (api *AdminApi) LoadFileList(c *gin.Context) {
	ctx := c.Request.Context()
	pageNo := c.Request.PostFormValue("pageNo")
	pageSize := c.Request.PostFormValue("pageSize")
	fileNameFuzzy := c.PostForm("fileNameFuzzy")
	filePid := c.PostForm("filePid")
	if pageNo == "" && pageSize == "" {
		pageNo = "1"
		pageSize = "20"
	}

	PageNo, err := strconv.ParseInt(pageNo, 10, 64)
	if err != nil {
		ginx.ResFail(c)
	}
	PageSize, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		ginx.ResFail(c)
	}

	res, err := api.AdminSrv.LoadFileList(ctx, int(PageNo), int(PageSize), fileNameFuzzy, filePid)
	if err != nil {
		ginx.ResFail(c)
	}

	ginx.ResOkWithData(c, res)
}

func (api *AdminApi) GetFolderInfo(c *gin.Context) {
	ctx := c.Request.Context()
	path := c.PostForm("path")

	srv := service.FileSrv{Repo: api.AdminSrv.FileRepo}

	res, err := srv.GetFolderInfo(ctx, path, contextx.FromUserID(ctx))

	if err != nil {
		ginx.ResFail(c)
		return
	}
	fmt.Printf("res.List: %v\n", res)
	ginx.ResOkWithData(c, res)

}

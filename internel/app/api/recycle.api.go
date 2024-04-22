package api

import (
	"fmt"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecycleApi struct {
	RecycleSrv *service.RecycleSrv
}

func (f *RecycleApi) GetFileList(c *gin.Context) {
	ctx := c.Request.Context()
	pageNo := c.DefaultPostForm("pageNo", "1")
	pageSize := c.DefaultPostForm("pageSize", "10")
	if pageNo == "" {
		pageNo = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}
	PageNo, err := strconv.ParseInt(pageNo, 10, 64)
	PageSize, err := strconv.ParseInt(pageSize, 10, 64)
	fmt.Printf("pageNo: %v\n", pageNo)
	fmt.Printf("pageNo: %v\n", pageSize)
	if err != nil {
		ginx.ResFailWithMessage(c, "获取数据失败")
		return
	}
	res, err := f.RecycleSrv.LoadListFiles(c, contextx.FromUserID(ctx), PageNo, PageSize)
	if err != nil {
		ginx.ResFailWithMessage(c, "获取数据失败")
		return
	}

	ginx.ResOkWithData(c, res)
}
func (f *RecycleApi) DelFiles(c *gin.Context) {
	ctx := c.Request.Context()

	input := c.PostForm("fileIds")
	if input == "" {
		ginx.ResFailWithMessage(c, "请选择文件夹")
		return
	}

	err := f.RecycleSrv.DelFiles(ctx, contextx.FromUserID(ctx), input)
	if err != nil {
		ginx.ResFailWithMessage(c, "删除失败")
		return
	}

	// f.RecycleSrv.Timer

	ginx.ResOk(c)
}

func (f *RecycleApi) RecoverFile(c *gin.Context) {
	ctx := c.Request.Context()

	fileIds := c.PostForm("fileIds")

	if err := f.RecycleSrv.RecoverFile(ctx, contextx.FromUserID(ctx), fileIds); err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOk(c)
}

package api

import (
	"fmt"
	"onlineCLoud/internel/app/dao/share"
	"onlineCLoud/internel/app/define"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"onlineCLoud/pkg/util/random"
	"onlineCLoud/pkg/util/uuid"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ShareApi struct {
	ShareSrv *service.ShareSrv
}

func (s *ShareApi) LoadShareList(c *gin.Context) {
	ctx := c.Request.Context()

	pageNo := c.PostForm("pageNo")
	pageSize := c.PostForm("pageSize")

	if pageNo == "" {
		pageNo = "1"
		pageSize = "20"
	}
	PageNo, err := strconv.ParseInt(pageNo, 10, 64)
	PageSize, err := strconv.ParseInt(pageSize, 10, 64)

	if err != nil {
		ginx.ResFail(c)
		return
	}

	list, err := s.ShareSrv.LoadShareList(ctx, contextx.FromUserID(ctx), PageNo, PageSize)
	if err != nil {
		ginx.ResFail(c)
		return
	}
	ginx.ResOkWithData(c, list)
}

func (s *ShareApi) ShareFile(c *gin.Context) {
	ctx := c.Request.Context()
	fileId := c.PostForm("fileId")
	validType := c.PostForm("validType")
	code := c.PostForm("code")
	validtype, err := strconv.ParseInt(validType, 10, 8)
	if fileId == "" || validType == "" || err != nil || validtype > 3 {
		ginx.ResFail(c)
		return
	}

	var share share.Share

	share.UserId = contextx.FromUserID(ctx)
	share.FileId = fileId
	share.ValidType = int8(validtype)
	share.ShareTime = time.Now().Format("2006-01-02 15:04:05")
	if validtype != define.FileShareForverDay {
		AddDay := 24 * define.GetDay(int8(validtype))
		share.ExpireTime = time.Now().Add(time.Hour * time.Duration(AddDay)).Format("2006-01-02 15:04:05")
	}

	if code == "" {
		code = random.GetStrRandom(5)
	}

	share.Code = code
	share.ShareId = uuid.MustString()

	err = s.ShareSrv.ShareFile(ctx, share)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, share)
}

// /share/cancelShare
func (s *ShareApi) CancelShare(c *gin.Context) {
	ctx := c.Request.Context()

	shareIds := c.PostFormArray("shareIds")

	err := s.ShareSrv.CancelShare(ctx, contextx.FromUserID(ctx), shareIds)

	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOk(c)
}

func (api *ShareApi) GetShareLoginInfo(c *gin.Context) {
	ctx := c.Request.Context()

	shareId := c.PostForm("shareId")

	info, err := api.ShareSrv.GetShareLoginInfo(ctx, shareId)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, info)
}
func (api *ShareApi) GetShareInfo(c *gin.Context) {
	ctx := c.Request.Context()

	shareId := c.PostForm("shareId")

	info, err := api.ShareSrv.GetShareInfo(ctx, shareId)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, info)
}

func (api *ShareApi) LoadFileList(c *gin.Context) {
	ctx := c.Request.Context()

	item := new(schema.RequestShareListPage)
	if err := ginx.ParseForm(c, item); err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	fmt.Println(item)
	info, err := api.ShareSrv.GetShareList(ctx, item)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, info)
}

func (api *ShareApi) GetFolderInfo(c *gin.Context) {
	ctx := c.Request.Context()

	shareId := c.PostForm("shareId")
	path := c.PostForm("path")

	info, err := api.ShareSrv.GetFolderInfo(ctx, shareId, path)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, info)
}

func (api *ShareApi) CheckShareCode(c *gin.Context) {
	ctx := c.Request.Context()

	shareId := c.PostForm("shareId")
	code := c.PostForm("code")

	info, err := api.ShareSrv.CheckShareCode(ctx, shareId, code)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	if info == true {
		ginx.ResOk(c)
	} else {
		ginx.ResFailWithMessage(c, "验证码错误")
	}

}

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"onlineCLoud/internel/app/dao/dto"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"onlineCLoud/pkg/util/random"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const sessionKey = "webShare_Key"
const StandTime = "2006-01-02 15:04:05"

type WebShareApi struct {
	ShareSrv *service.ShareSrv
	FileSrv  *service.FileSrv
}

func (api *WebShareApi) GetShareLoginInfo(c *gin.Context) {
	ctx := c.Request.Context()
	shareId := c.PostForm("shareId")
	if shareId == "" {
		ginx.ResFailWithMessage(c, "请求内容为空")
		return
	}

	if err := api.checkShare(ctx, shareId); err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	session := sessions.Default(c)
	if share := session.Get(sessionKey + shareId); share == nil {
		ginx.ResOk(c)
		return
	}

	info, err := api.ShareSrv.GetShareInfo(ctx, shareId)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, info)
}

func (api *WebShareApi) GetShareInfo(c *gin.Context) {
	ctx := c.Request.Context()

	shareId := c.PostForm("shareId")

	info, err := api.ShareSrv.GetShareInfo(ctx, shareId)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, info)
}

func (api *WebShareApi) LoadFileList(c *gin.Context) {

	ctx := c.Request.Context()
	p1 := time.Now()
	item := new(schema.RequestShareListPage)
	if err := ginx.ParseForm(c, item); err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	p2 := time.Now()

	session := dto.GetSession(c, item.ShareId)
	if session == nil {
		return
	}

	p3 := time.Now()
	if err := api.checkShare(ctx, item.ShareId); err != nil {
		ginx.ResOkWithMessage(c, err.Error())
		return
	}

	p4 := time.Now()
	info, err := api.ShareSrv.GetShareList(ctx, item)
	if err != nil {
		ginx.ResFail(c)
		return
	}
	p5 := time.Now()
	ginx.ResOkWithData(c, info)

	fmt.Printf("p2.Sub(p1).Milliseconds(): %v\n", p2.Sub(p1).Milliseconds())
	fmt.Printf("p3.Sub(p2).Milliseconds(): %v\n", p3.Sub(p2).Milliseconds())
	fmt.Printf("p2.Sub(p1).Milliseconds(): %v\n", p4.Sub(p3).Milliseconds())
	fmt.Printf("p2.Sub(p1).Milliseconds(): %v\n", p5.Sub(p4).Milliseconds())

}

func (api *WebShareApi) GetFolderInfo(c *gin.Context) {
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

func (api *WebShareApi) CheckShareCode(c *gin.Context) {
	ctx := c.Request.Context()

	shareId := c.PostForm("shareId")
	code := c.PostForm("code")
	if shareId == "" || code == "" {
		ginx.ResOkWithMessage(c, "参数错误")
		return
	}
	info, err := api.ShareSrv.CheckShareCode(ctx, shareId, code)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	if info == nil {
		ginx.ResOkWithMessage(c, "验证码错误")
	} else {
		session := sessions.Default(c)
		bodyv, _ := json.Marshal(info)
		session.Set(sessionKey+shareId, bodyv)
		session.Save()
		ginx.ResOk(c)
	}

}
func (api *WebShareApi) checkShare(ctx context.Context, shareId string) error {

	info, err := api.ShareSrv.GetShareInfo(ctx, shareId)
	if err != nil {
		return err
	}
	if info == nil {
		return errors.New("分享信息不存在")
	}
	if info.ExpireTime != "永久" {
		now := time.Now()

		ex, err := time.Parse(StandTime, info.ExpireTime)
		if err != nil {
			return err
		}

		if now.After(ex) {
			return errors.New("分享的信息失效")
		}

	}
	return nil
}

func (api *WebShareApi) GetFile(c *gin.Context) {

	ctx := c.Request.Context()
	shareId := c.Param("shareId")
	fileId := c.Param("fileId")
	if len(shareId) == 0 || len(fileId) == 0 {
		ginx.ResData(c, http.StatusBadRequest, []byte("参数错误"))
		return
	}
	session := dto.GetSession(c, shareId)
	if session == nil {
		return
	}
	body, err := api.FileSrv.GetFile(ctx, fileId, session.ShareUserId)
	if err != nil {
		ginx.ResData(c, 500, []byte(err.Error()))
		return
	}

	ginx.ResData(c, 200, body)
}
func (api *WebShareApi) CreateDownloadUrl(c *gin.Context) {

	ctx := c.Request.Context()
	shareId := c.Param("shareId")
	fileId := c.Param("fileId")
	if len(shareId) == 0 || len(fileId) == 0 {
		ginx.ResData(c, http.StatusBadRequest, []byte("参数错误"))
		return
	}

	session := dto.GetSession(c, shareId)
	if session == nil {
		return
	}

	file, err := api.FileSrv.GetFileInfo(ctx, fileId, session.ShareUserId)
	if err != nil || file == nil || file.CreateTime == "" || file.FolderType == 1 {
		ginx.ResJson(c, 600, "", "操作错误", "fail")
		return
	}
	code := random.GetStrRandom(50)

	dto := dto.DownloadDto{
		Code:     code,
		FileName: file.FileName,
		Path:     file.FilePath,
	}
	v, _ := json.Marshal(dto)
	rdx := redisx.NewClient()
	rdx.Set(ctx, fmt.Sprintf("download:%v", code), string(v), time.Duration(30)*time.Minute)
	ginx.ResOkWithData(c, code)

}

func (api *WebShareApi) Download(c *gin.Context) {

	ctx := c.Request.Context()
	code := c.Param("code")

	if len(code) == 0 {
		ginx.ResFail(c)
		return
	}

	var dto dto.DownloadDto
	rdx := redisx.NewClient()
	str, err := rdx.Get(ctx, fmt.Sprintf("download:%v", code))
	if err != nil {
		ginx.ResFail(c)
		return
	}

	json.Unmarshal([]byte(str), &dto)

	file, err := os.Open(dto.Path)
	if err != nil {
		c.String(http.StatusNotFound, "文件未找到")
		return
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "读取文件信息时出错")
		return
	}
	// 限制每次读取的字节数为 100 KB
	limitedReader := &RateLimitedReader{
		R:     file,
		Limit: 1024 * 1024, // 每秒读取100 KB
	}
	c.Header("Content-Length", fmt.Sprintf("%d", fi.Size()))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", dto.FileName))
	http.ServeContent(c.Writer, c.Request, dto.FileName, fi.ModTime(), limitedReader)
}

func (api *WebShareApi) SaveShare(c *gin.Context) {
	ctx := c.Request.Context()
	shareId := c.PostForm("shareId")
	shareFileIds := c.PostForm("shareFileIds")
	myFolderId := c.PostForm("myFolderId")

	if len(shareId) == 0 || len(shareFileIds) == 0 || len(myFolderId) == 0 {
		ginx.ResData(c, 400, []byte("参数缺失"))
		return
	}

	session := dto.GetSession(c, shareId)
	if session == nil {
		ginx.ResData(c, 403, []byte("数据获取错误"))
		return
	}

	currentUser := contextx.FromUserID(ctx)

	if currentUser == session.ShareUserId {
		ginx.ResData(c, 403, []byte("自己分享的文件无法保存到自己"))
		return
	}

	err := api.FileSrv.SaveShare(ctx, session.FileId, shareFileIds, myFolderId, session.ShareUserId, currentUser)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOk(c)
}
func (api *WebShareApi) GetVideoInfo(c *gin.Context) {
	ctx := c.Request.Context()
	shareId := c.Param("shareId")
	fid := c.Param("fid")

	if len(shareId) == 0 || len(fid) == 0 {
		ginx.ResData(c, 400, []byte("参数缺失"))
		return
	}

	session := dto.GetSession(c, shareId)
	if session == nil {
		ginx.ResData(c, 403, nil)
		return
	}

	body, err := api.FileSrv.GetFile(ctx, fid, session.ShareUserId)

	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResData(c, 200, []byte(body))
}

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"onlineCLoud/internel/app/dao/dto"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const sessionKey = "webShare_Key"
const StandTime = "2006-01-02 15:04:05"

type WebShareApi struct {
	UserSrv     *service.UserSrv
	ShareSrv    *service.ShareSrv
	FileSrv     *service.FileSrv
	DownLoadSrv *service.DownLoadSrv
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
	item := new(schema.RequestShareListPage)
	if err := ginx.ParseForm(c, item); err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	session := dto.GetSession(c, item.ShareId)
	if session == nil {
		ginx.ResNeedReload(c)
		return
	}

	if err := api.checkShare(ctx, item.ShareId); err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	info, err := api.ShareSrv.GetShareList(ctx, item)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, info)

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
		ginx.ResNeedReload(c)
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
		ginx.ResNeedReload(c)
		return
	}

	file, err := api.FileSrv.GetFileInfo(ctx, fileId, session.ShareUserId)
	if err != nil || file == nil || file.CreateTime.IsZero() || file.FolderType == 1 {
		ginx.ResJson(c, 600, "", "操作错误", "fail")
		return
	}

	currentime := time.Now().Unix()
	code := fmt.Sprintf("%v_%v_%v", file.UserID, currentime, file.FileID)

	err = api.DownLoadSrv.CreateDownLoad(ctx, file.UserID, file.FileID, file.FilePath, code)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOkWithData(c, code)

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

	// 文件防护   ----->>
	fileIds := strings.Split(shareFileIds, ",")
	for _, id := range fileIds {
		err := api.FileSrv.CheckFootFilePid(ctx, session.FileId, session.ShareUserId, id)
		if err != nil {
			ginx.ResFailWithMessage(c, err.Error())
			return
		}
	}

	// 检查剩余空间

	sum, err := api.FileSrv.GetFileListTotalSize(ctx, session.ShareUserId, fileIds)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	space := api.UserSrv.GetUserSpaceById(ctx, currentUser)
	if space.TotalSpace == 0 {
		ginx.ResFailWithMessage(c, errors.New("获取错误").Error())
		return
	}

	if space.UseSpace+sum > space.TotalSpace {
		ginx.ResFailWithMessage(c, errors.New("空间不够").Error())
		return
	}

	err = api.FileSrv.SaveShare(ctx, session.FileId, shareFileIds, myFolderId, session.ShareUserId, currentUser)
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
		ginx.ResNeedReload(c)
		return
	}

	body, err := api.FileSrv.GetFile(ctx, fid, session.ShareUserId)

	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResData(c, 200, []byte(body))
}

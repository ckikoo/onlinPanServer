package api

import (
	"fmt"
	"log"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"

	"time"

	"github.com/gin-gonic/gin"
)

type FileApi struct {
	FileSrv     *service.FileSrv
	DownLoadSrv *service.DownLoadSrv
}

func (f *FileApi) GetFileList(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.RequestFileListPage
	if item.PageNo == 0 || item.PageSize == 0 {
		item.PageNo = 1
		item.PageSize = 20
	}
	err := ginx.ParseForm(c, &item)
	if err != nil {
		log.Default().Printf("请求继续失败,接口GETFILELIST %v\n", err)
		ginx.ResFailWithMessage(c, "数据格式有误")
		return
	}
	res, err := f.FileSrv.LoadListFiles(c, contextx.FromUserID(ctx), &item)
	if err != nil {
		ginx.ResFailWithMessage(c, "获取数据失败")
		return
	}

	ginx.ResOkWithData(c, res)
}

func (f *FileApi) UploadFile(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.FileUpload
	if err := ginx.ParseForm(c, &item); err != nil {
		ginx.ResFailWithMessage(c, "数据格式有误")
		return
	}

	op, err := f.FileSrv.UploadFile(c, contextx.FromUserID(ctx), item)
	if err != nil {
		ginx.ResFailWithMessage(c, "上传失败")
		return
	}
	ginx.ResOkWithData(c, op)

}

func (f *FileApi) CancelUpload(c *gin.Context) {
	ctx := c.Request.Context()
	fileId := c.PostForm("fileId")
	m := make(map[string]interface{}, 0)
	item := f.FileSrv.Timer.Del(fileId + contextx.FromUserID(ctx))
	item.Action()
	m["status"] = "OK"

	ginx.ResOkWithData(c, m)
}

func (f *FileApi) NewFoloder(c *gin.Context) {
	ctx := c.Request.Context()

	filePid := c.PostForm("filePid")
	fileName := c.PostForm("fileName")

	info, err := f.FileSrv.NewFoloder(c, contextx.FromUserID(ctx), filePid, fileName)
	if err != nil {
		ginx.ResFailWithMessage(c, "创建失败")
		return
	}
	ginx.ResOkWithData(c, info)

}

func (f *FileApi) DelFiles(c *gin.Context) {
	ctx := c.Request.Context()

	input := c.PostForm("fileIds")
	if input == "" {
		ginx.ResFailWithMessage(c, "请选择文件夹")
		return
	}

	err := f.FileSrv.DelFiles(c, contextx.FromUserID(ctx), input)
	if err != nil {
		ginx.ResFailWithMessage(c, "删除失败")
		return
	}
	ginx.ResOk(c)
}

func (f *FileApi) GetImage(c *gin.Context) {
	imgsrc := c.Param("src")

	f.FileSrv.GetImage(c.Writer, c.Request, imgsrc)

}
func (f *FileApi) GetVideoInfo(c *gin.Context) {
	ctx := c.Request.Context()
	fid := c.Param("fid")
	if fid == "" {
		ginx.ResFail(c)
	}

	body, err := f.FileSrv.GetFile(ctx, fid, contextx.FromUserID(ctx))

	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResData(c, 200, []byte(body))
}

func (f *FileApi) GetFileInfo(c *gin.Context) {
	ctx := c.Request.Context()
	fid := c.Param("fid")
	if fid == "" {
		ginx.ResFail(c)
	}
	body, err := f.FileSrv.GetFile(ctx, fid, contextx.FromUserID(ctx))

	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResData(c, 200, []byte(body))
}

func (f *FileApi) GetFolderInfo(c *gin.Context) {
	ctx := c.Request.Context()
	path := c.PostForm("path")

	res, err := f.FileSrv.GetFolderInfo(ctx, path, contextx.FromUserID(ctx))

	if err != nil {
		ginx.ResFail(c)
		return
	}
	ginx.ResOkWithData(c, res)

}

func (f *FileApi) FileRename(c *gin.Context) {
	ctx := c.Request.Context()
	fileId := c.PostForm("fileId")
	filePId := c.PostForm("filePid")
	fileName := c.PostForm("fileName")

	if fileName == "" {
		ginx.ResFailWithMessage(c, "文件名不能为空")
		return
	}
	if fileId == "" || filePId == "" {
		log.Println("用户参数不合法")
		ginx.ResFail(c)
		return
	}
	file, err := f.FileSrv.FileRename(ctx, contextx.FromUserID(ctx), fileId, filePId, fileName)
	if err != nil || file == nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, file)

}

func (f *FileApi) LoadAllFolder(c *gin.Context) {
	ctx := c.Request.Context()
	filePid := c.PostForm("filePid")
	currentFileIds := c.PostForm("currentFileIds")
	if filePid == "" {
		log.Println("LoadAllFolder 用户参数不合法")
		ginx.ResFail(c)
	}
	files, err := f.FileSrv.LoadAllFolder(ctx, contextx.FromUserID(ctx), filePid, currentFileIds)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, files)

}

func (f *FileApi) ChangeFileFolder(c *gin.Context) {
	ctx := c.Request.Context()

	fileIds := c.Request.FormValue("fileIds")
	filePid := c.Request.FormValue("filePid")
	if filePid == "" || fileIds == "" {
		log.Println("LoadAllFolder 用户参数不合法")
		ginx.ResFail(c)
	}
	fmt.Println(filePid)
	err := f.FileSrv.ChangeFileFolder(ctx, contextx.FromUserID(ctx), fileIds, filePid)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	ginx.ResOk(c)
}

func (f *FileApi) CreateDownloadUrl(c *gin.Context) {
	ctx := c.Request.Context()

	fid := c.Param("fid")
	file, err := f.FileSrv.GetFileInfo(ctx, fid, contextx.FromUserID(ctx))
	if err != nil || file == nil || file.CreateTime == "" || file.FolderType == 1 {
		ginx.ResJson(c, 400, "", "操作错误", "fail")
		return
	}
	currentime := time.Now().Unix()

	code := fmt.Sprintf("%v_%v_%v", file.UserID, currentime, file.FileID)

	err = f.DownLoadSrv.CreateDownLoad(ctx, file.UserID, file.FileID, file.FilePath, code)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	ginx.ResOkWithData(c, code)
}

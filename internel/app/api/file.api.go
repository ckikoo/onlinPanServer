package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"onlineCLoud/internel/app/dao/dto"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"onlineCLoud/pkg/util/json"
	"onlineCLoud/pkg/util/random"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type FileApi struct {
	FileSrv *service.FileSrv
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

	err := os.RemoveAll(fmt.Sprintf("temp/%v/%v/%v", time.Now().Month(), contextx.FromUserID(ctx), fileId))
	if err != nil {
		m["status"] = "ERROR"
	} else {
		m["status"] = "OK"
	}

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
	fmt.Println("debug")
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
	fmt.Printf("res.List: %v\n", res)
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
	fmt.Printf("contextx.FromUserID(ctx): %v\n", contextx.FromUserID(ctx))
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
	fmt.Printf("contextx.FromUserID(ctx): %v\n", contextx.FromUserID(ctx))
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
		panic(err)
		return
	}
	ginx.ResOk(c)
}

func (f *FileApi) CreateDownloadUrl(c *gin.Context) {
	ctx := c.Request.Context()

	fid := c.Param("fid")
	file, err := f.FileSrv.GetFileInfo(ctx, fid, contextx.FromUserID(ctx))
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

	rdx := redisx.NewClient()
	rdx.Set(ctx, fmt.Sprintf("download:%v", code), json.MarshalToString(dto), time.Duration(30)*time.Minute)
	ginx.ResOkWithData(c, code)
}

func (f *FileApi) Download(c *gin.Context) {
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

// RateLimitedReader 实现了 io.Reader 接口，用于限制读取速度。
type RateLimitedReader struct {
	R       io.Reader // 原始的 io.Reader
	Limit   int64     // 每秒读取的字节数限制
	LastSec int64     // 上次读取的时间戳
	ReadCnt int64     // 当前秒内已读取的字节数
}

func (r *RateLimitedReader) Seek(offset int64, whence int) (int64, error) {
	return r.R.(*os.File).Seek(offset, whence)
}

func (r *RateLimitedReader) Read(p []byte) (n int, err error) {
	// 获取当前时间戳
	now := time.Now().Unix()

	// 如果距离上次读取的时间大于1秒，则重置读取字节数
	if now-r.LastSec > 1 {
		r.ReadCnt = 0
		r.LastSec = now
	}

	// 计算当前秒内还可以读取的字节数
	remaining := r.Limit - r.ReadCnt

	// 如果剩余可读字节数为0，则等待1秒后重新计算
	if remaining <= 0 {
		time.Sleep(400 * time.Millisecond)
		return 0, nil
	}

	// 限制每次读取的字节数不超过剩余可读字节数
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}

	// 读取数据
	n, err = r.R.Read(p)
	if err != nil {
		return n, err
	}

	// 更新已读取字节数
	r.ReadCnt += int64(n)
	return n, nil
}

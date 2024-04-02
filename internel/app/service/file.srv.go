package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"onlineCLoud/internel/app/dao/file"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/define"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/pkg/contextx"
	fileUtil "onlineCLoud/pkg/util/file"
	processutil "onlineCLoud/pkg/util/process"
	"onlineCLoud/pkg/util/uuid"
	util "onlineCLoud/pkg/util/uuid"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FileSrv struct {
	Repo *file.FileRepo
}

func (f *FileSrv) LoadListFiles(ctx context.Context, uid string, p *schema.RequestFileListPage) (*schema.ListResult, error) {

	var res schema.ListResult
	p.DelFlag = define.FileFlagInUse
	if p.Category != "" && p.Category != "all" {
		p.Category = strconv.Itoa(int(define.FileCategoryStr4ID(p.Category)))
	}
	files, err := f.Repo.GetFileList(ctx, uid, p, true)
	if err != nil {
		return nil, err
	}
	total, err := f.Repo.GetFileListTotal(ctx, uid, p)
	if err != nil {
		return nil, err
	}

	res.List = files
	res.Parms = &p.PageParams
	res.PageTotal = (total + int64(p.GetPageSize())/2) / int64(p.GetPageSize())
	res.TotalCount = total
	return &res, nil
}

func (srv *FileSrv) UploadFilePre(c context.Context, email string, fileSize uint64) (map[string]any, error) {
	resMap := make(map[string]interface{}, 0)

	userdb := user.GetUserDB(c, srv.Repo.Db)
	usrv := UserSrv{UserRepo: &user.UserRepo{
		DB: userdb,
		Rd: redisx.NewClient(),
	}}
	spaceMap := usrv.GetUserSpace(c, email)
	if spaceMap == nil {
		resMap["status"] = "INTERNET_ERROR"
		return resMap, errors.New("no found info")
	}

	var space user.UserSpace
	json.Unmarshal([]byte(fmt.Sprintf("%v", spaceMap)), &space)
	if space.UseSpace+fileSize > space.TotalSpace {
		resMap["status"] = Uer_NO_SPACE
	} else {
		resMap["status"] = Uer_SPACE_SA
	}

	return resMap, nil
}

// 返回状态集， error
func (srv *FileSrv) UploadFile(c *gin.Context, uid string, fileinfo schema.FileUpload) (map[string]interface{}, error) {
	statusMap := make(map[string]interface{}, 0)
	usrv := UserSrv{UserRepo: &user.UserRepo{DB: user.GetUserDB(c.Request.Context(), srv.Repo.Db), Rd: redisx.NewClient()}}
	userspace := usrv.GetUserSpace(c.Request.Context(), contextx.FromUserEmail(c.Request.Context()))

	s, _ := json.Marshal(userspace)

	var space user.UserSpace
	json.Unmarshal([]byte(s), &space)

	if space.UseSpace+uint64(fileinfo.FileSize) > space.TotalSpace {
		statusMap["status"] = "fail"
		statusMap["errorMsg"] = "空间不足"
		return statusMap, nil
	}
	if fileinfo.ChunkIndex == 0 {
		fmt.Println("debug the first chunk")
		file, err := srv.Repo.CheckFileName(c.Request.Context(), fileinfo.FilePid, uid, fileinfo.FileName, "0")
		if err != nil {
			return nil, err
		}

		fmt.Printf("file: %v\n", file)
		if file != nil && file.FileName != "" {
			fmt.Println("111")
			fileinfo.FileName = fileUtil.Rename(fileinfo.FileName)
		}

		if file != nil && file.FileMd5 != "" {
			file.FileName = fileinfo.FileName
			file.CreateTime = time.Now().Format("2006-01-02 15:04:05")
			file.LastUpdateTime = time.Now().Format("2006-01-02 15:04:05")
			file.FileID = util.MustString()
			file.DelFlag = define.FileFlagInUse
			file.Status = define.FileStatusUsing // 成功过
			file.FilePid = fileinfo.FilePid
			file.RecoveryTime = ""
			file.UserID = uid
			if err := srv.Repo.UploadFile(c, file); err != nil {
				return nil, nil
			}
			usrv.UpdateSpace(c, contextx.FromUserEmail(c.Request.Context()), file.FileSize)

			statusMap["fileId"] = file.FileID
			statusMap["status"] = FILE_STATUS_USING
			return statusMap, nil
		}
	}

	// 第一片分片文件不存在+ 其他完整了
	fh, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	src, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	buf := make([]byte, min(50*1024*1024, fh.Size))
	// 文件夹路劲 upload/ month / 用户名 / fileid
	var fileid string
	if fileinfo.FileId == "" {
		fileid = uuid.MustString()
	} else {
		fileid = fileinfo.FileId
	}
	tempDir := fmt.Sprintf("temp/%v/%v/%v", time.Now().Month(), uid, fileid)
	filePath := fmt.Sprintf("%v/%v", tempDir, fileinfo.ChunkIndex)
	newFile, err := fileUtil.FileCreate(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, errors.New("")
	}
	defer newFile.Close()
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return nil, errors.ErrUnsupported
		}
		if n == 0 {
			break
		}
		_, err = newFile.Write(buf)
		if err != nil {
			log.Fatal("空间不足")
			return nil, errors.ErrUnsupported
		}
	}
	statusMap["fileId"] = fileid
	statusMap["status"] = FILE_STATUS_TRANSFER
	if fileinfo.ChunkIndex == fileinfo.Chunks-1 {
		// upload/FileId/FileName
		uploadDir := fmt.Sprintf("upload/%v", fileid)
		uploadFile := fmt.Sprintf("%v/%v", uploadDir, fileinfo.FileName)
		if err := fileUtil.FileMerge(tempDir, uploadFile); err != nil {
			statusMap["fileId"] = fileid
			statusMap["status"] = FILE_STATUS_TRANSFER_FAIL
			return statusMap, nil
		}
		var file file.File
		file.FileID = fileid
		file.UserID = uid
		file.FileName = fileinfo.FileName
		file.FileSize = uint64(fileinfo.FileSize)
		file.FileMd5 = fileinfo.FileMd5
		file.FilePath = uploadFile
		file.DelFlag = define.FileFlagInUse
		file.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		file.LastUpdateTime = time.Now().Format("2006-01-02 15:04:05")
		file.Status = define.FileStatusUsing
		file.FilePid = fileinfo.FilePid
		file.FolderType = 0
		ext := fileUtil.GetFileExt(file.FileName)
		file.FileType = define.GetFileType(ext)
		file.FileCategory = define.FileCategoryStr4ID(ext)
		usrv.UpdateSpace(c.Request.Context(), contextx.FromUserEmail(c.Request.Context()), file.FileSize)
		//TODO
		if file.FileType == define.FileTypeVideo {
			CutFile4Video(fileid, file.FilePath) // 文件切片
			dest := fmt.Sprintf("%s/%s", uploadDir, file.FileID+".png")
			CreateCover4Video(file.FilePath, 150, dest)
			file.FileCover = fmt.Sprintf("%v", file.FileID+".png")
		} else if file.FileType == define.FileTypeImage {
			//生成缩略图
			dest := fmt.Sprintf("%s/%s", uploadDir, file.FileID+".png")
			CreateCover4Video(file.FilePath, 150, dest)
			file.FileCover = fmt.Sprintf("%v", file.FileID+".png")
		}
		if err := srv.Repo.UploadFile(c, &file); err != nil {
			fmt.Println("upload file error", err)
			os.RemoveAll(file.FilePath)
			statusMap["fileId"] = fileid
			statusMap["status"] = FILE_STATUS_TRANSFER_FAIL
			return statusMap, nil
		}
		statusMap["fileId"] = fileid
		statusMap["status"] = FILE_STATUS_USING

		err := os.RemoveAll(tempDir)
		if err != nil {
			statusMap["satus"] = FILE_STATUS_TRANSFER_FAIL
		}
	}

	return statusMap, nil
}

func CutFile4Video(fileId, videoFilePath string) error {
	path, err := fileUtil.NewDir(videoFilePath[:strings.LastIndex(videoFilePath, "/")])
	if err != nil {
		log.Printf("cutfile4video error creating video folder: %s", err)
		return err
	}

	tsPath := fmt.Sprintf("%v/%v%v", path, fileId, ".ts")
	cmd := exec.Command("ffmpeg", "-y", "-i", videoFilePath, "-vcodec", "copy", "-acodec", "copy", "-bsf:v", "h264_mp4toannexb", tsPath)

	// 指定命令生成 .ts 文件
	if _, err = processutil.ExecuteCommand(cmd, false); err != nil {
		log.Printf("cutfile4video error during .ts file generation: %v", err)
		return err
	}

	// 生成索引文件（m3u8）并进行切片
	cmd = exec.Command("ffmpeg", "-y", "-i", videoFilePath, "-vcodec", "copy", "-acodec", "copy", "-bsf:v", "h264_mp4toannexb", "-f", "hls", "-hls_time", "30", "-hls_list_size", "0", "-hls_segment_filename", fmt.Sprintf("%v/%v_%%d.ts", path, fileId), path+"/index.m3u8")

	// 分片
	if _, err := processutil.ExecuteCommand(cmd, false); err != nil {
		log.Printf("cutfile4video error during .m3u8 and ts file generation: %v", err)
		return err
	}

	// 删除 index.ts 文件
	if err := os.Remove(tsPath); err != nil {
		log.Printf("cutfile4video error deleting index.ts: %v", err)
	}

	return nil
}

func (f *FileSrv) CheckFileName(ctx context.Context, filePid string, userID string, fileName string, folderType string) (*file.File, error) {
	file, err := f.Repo.CheckFileName(ctx, filePid, userID, fileName, folderType)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *FileSrv) NewFoloder(ctx context.Context, uid string, filePid, fileName string) (*file.File, error) {

	tmp, err := f.CheckFileName(ctx, filePid, uid, fileName, "1")
	if err != nil {
		return nil, err
	}
	if tmp != nil && tmp.FileID != "" {
		return nil, errors.New("文件已经存在")
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	file := file.File{FileID: uuid.MustString(), UserID: uid, FolderType: 1, FileName: fileName, FilePid: filePid, CreateTime: now, LastUpdateTime: now, RecoveryTime: now, Status: define.FileStatusUsing, DelFlag: define.FileFlagInUse}
	err = f.Repo.UploadFile(ctx, &file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (f *FileSrv) findAllSubFolderFileList(ctx context.Context, fileIdList *[]string, userID, fileID string, delflag int8) {
	*fileIdList = append(*fileIdList, fileID)

	query := schema.RequestFileListPage{
		FilePid:    fileID,
		DelFlag:    delflag,
		FolderType: 1,
	}

	fields, err := f.Repo.GetFileList(ctx, userID, &query, false)

	if err != nil || fields == nil || len(fields) == 0 {
		return
	}

	for _, v := range fields {
		f.findAllSubFolderFileList(ctx, fileIdList, userID, v.FileID, delflag)
	}

	fmt.Printf("level fileIdList: %v\n", fileIdList)
}

func (f *FileSrv) DelFiles(ctx context.Context, uid string, fileId string) error {
	fileIds := strings.Split(fileId, ",")
	query := schema.RequestFileListPage{
		Path:    fileIds,
		DelFlag: define.FileFlagInUse,
	}

	// 查找目录
	fileInfoList, _ := f.Repo.GetFileList(ctx, uid, &query, false)
	if fileInfoList == nil || len(fileInfoList) == 0 {
		return nil
	}
	delFileList := make([]string, 0)
	//TODO fix

	for _, e := range fileInfoList {
		f.findAllSubFolderFileList(ctx, &delFileList, uid, e.FileID, define.FileFlagInUse)
	}

	// 暂时移除下面的子目录
	if len(delFileList) != 0 {
		err := f.Repo.UpdateFileDelFlag(ctx, uid, delFileList, nil, define.FileFlagInUse, define.FileFlagSoftDeleted, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Default().Printf("update file delFileList error:%v", err)
			return err
		}
	}
	err := f.Repo.UpdateFileDelFlag(ctx, uid, nil, fileIds, define.FileFlagInUse, define.FileFlagInRecycleBin, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Default().Printf("update file delflag error:%v", err)
		return err
	}

	return nil
}

func CreateCover4Video(path string, width int, desc string) error {
	cmd := exec.Command("/usr/bin/ffmpeg", "-i", path, "-y", "-vframes", "1", "-vf", fmt.Sprintf("scale=%d:%d", width, width), desc)
	cmd.Stderr = os.Stderr
	if _, err := processutil.ExecuteCommand(cmd, false); err != nil {
		log.Printf("create cover error: %v", err)
		return err
	}
	return nil
}

func (f *FileSrv) GetImage(w http.ResponseWriter, r *http.Request, name string) {
	fid := name[:strings.LastIndex(name, ".")]

	var content bytes.Buffer

	file, err := os.Open(fmt.Sprintf("upload/%s/%s", fid, name))
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer file.Close()

	// 文件读取到
	_, err = io.Copy(&content, file)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Cache-Control", "max-age=2500")
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(content.Bytes()))
}

func (f *FileSrv) GetFile(ctx context.Context, fid string, uid string) ([]byte, error) {
	var flag bool = false
	tmp := ""
	if strings.HasSuffix(fid, ".ts") {
		flag = true
		tmp = fid
		fid = fid[:strings.LastIndex(fid, "_")]
	}

	file, err := f.Repo.GetFileInfo(ctx, fid, uid)
	if err != nil {
		return make([]byte, 0), err
	}

	if nil == file {
		return make([]byte, 0), nil
	}
	if flag {
		fileNameNoSuffix := fmt.Sprintf("%v/%v/%v", "upload", fid, tmp)
		b, err := os.ReadFile(fileNameNoSuffix)
		if err != nil {
			return make([]byte, 0), err
		}
		return b, nil
	} else {
		if (file.FileType) == define.FileTypeVideo {
			fileNameNoSuffix := fmt.Sprintf("%v%v%v", "upload/", fid, "/index.m3u8")
			b, err := os.ReadFile(fileNameNoSuffix)
			if err != nil {
				return make([]byte, 0), err
			}
			return b, nil

		} else {
			filePath := fmt.Sprintf("%v/%v/%v", "upload", fid, file.FileName)
			b, err := os.ReadFile(filePath)
			if err != nil {
				return make([]byte, 0), err
			}
			return b, nil
		}
	}
}

func (f *FileSrv) GetFolderInfo(ctx context.Context, path string, uid string) ([]file.File, error) {
	paths := strings.Split(path, "/")
	fmt.Printf("path=========: %v\n", path)
	var item schema.RequestFileListPage
	item.Path = paths
	item.FolderType = 1
	item.DelFlag = define.FileFlagInUse
	fmt.Printf("item: %v\n", item)
	res, err := f.Repo.GetFileList(ctx, uid, &item, false)
	if err != nil {
		return nil, err
	}
	fmt.Printf("item: %v\n", res)
	return res, nil

}

func (f *FileSrv) GetFileInfo(ctx context.Context, fid string, uid string) (*file.File, error) {

	res, err := f.Repo.GetFileInfo(ctx, fid, uid)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (f *FileSrv) FileRename(ctx context.Context, uid, fileID, filePid, fileName string) (*file.File, error) {

	file, err := f.Repo.GetFileInfo(ctx, fileID, uid)
	if err != nil {
		log.Printf("filerename error %v:", err)
		return nil, err
	}

	if file == nil || file.CreateTime == "" {
		return nil, err
	}

	if file.FolderType == 0 {
		fileName = fileName + "." + fileUtil.GetFileExt(file.FileName)
	}
	tmp, err := f.Repo.CheckFileName(ctx, filePid, uid, fileName, fmt.Sprintln(file.FolderType))
	if err != nil {
		log.Printf("filerename CheckFileName error %v:", err)
		return nil, err
	}

	if tmp != nil && tmp.CreateTime != "" {
		return nil, errors.New("文件存在")
	}

	file.FileName = fileName
	file.LastUpdateTime = time.Now().Format("2006-01-02 15:04:05")

	err = f.Repo.UpdateFile(ctx, file)
	if err != nil {
		log.Printf("filerename UploadFile error %v:", err)
		return nil, err
	}
	return file, nil
}

func (f *FileSrv) LoadAllFolder(ctx context.Context, uid string, filePid string, fileIDs string) ([]file.File, error) {
	var item schema.RequestFileListPage
	curs := strings.Split(fileIDs, ",")
	item.ExInclude = curs
	item.FolderType = 1
	item.DelFlag = define.FileFlagInUse
	item.FilePid = filePid
	res, err := f.Repo.GetFileList(ctx, uid, &item, false)
	return res, err
}
func (f *FileSrv) ChangeFileFolder(ctx context.Context, uid string, fileIds string, filePid string) ([]file.File, error) {

	if strings.Contains(fileIds, filePid) {
		return nil, errors.New("")
	}
	if filePid != "0" { //判定父文件夹是否存在
		file, err := f.Repo.GetFileInfo(ctx, filePid, uid)
		if err != nil || file == nil || file.Status != define.FileStatusUsing {
			return nil, errors.New("error") // TODO wait fix
		}
	}

	var item schema.RequestFileListPage

	item.FilePid = filePid
	lists, err := f.Repo.GetFileList(ctx, uid, &item, false) // 找新文件夹下的子文件
	if err != nil {
		log.Println("ChangeFileFolder GetFileList error", err)
		return nil, errors.New("error") // TODO wait fix
	}

	fileNameMap := make(map[string]file.File, 0)
	for _, v := range lists {
		fileNameMap[v.FileName] = v
	}

	item = schema.RequestFileListPage{}
	curs := strings.Split(fileIds, ",")
	item.Path = curs
	lists, err = f.Repo.GetFileList(ctx, uid, &item, false) // 找要被移动的文件

	for _, ele := range lists {
		if _, ok := fileNameMap[ele.FileName]; ok {
			fileNewName := fileUtil.Rename(ele.FileName)
			ele.FileName = fileNewName
		}
		ele.FilePid = filePid
		f.Repo.UpdateFile(ctx, &ele)
	}

	return nil, nil
}

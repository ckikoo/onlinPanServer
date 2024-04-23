package service

import (
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
	"onlineCLoud/pkg/cache"
	"onlineCLoud/pkg/contextx"
	"onlineCLoud/pkg/timer"
	fileUtil "onlineCLoud/pkg/util/file"
	hdfsUtil "onlineCLoud/pkg/util/hdfs"
	ossUtil "onlineCLoud/pkg/util/oss"
	processutil "onlineCLoud/pkg/util/process"
	"onlineCLoud/pkg/util/uuid"
	util "onlineCLoud/pkg/util/uuid"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
)

type FileSrv struct {
	Repo  *file.FileRepo
	Timer *timer.TimerManager
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
		file, err := srv.Repo.CheckFileName(c.Request.Context(), fileinfo.FilePid, uid, fileinfo.FileName, "0")
		if err != nil {
			return nil, err
		}
		fmt.Printf("file: %v\n", file)
		if file != nil && file.FileName != "" {
			file.FileName = fileUtil.Rename(fileinfo.FileName)
		}

		file, err = srv.Repo.GetFileByMd5(c, fileinfo.FileMd5)
		if err != nil {
			return nil, err
		}

		if file != nil && file.FileMd5 != "" {
			file.FileName = fileUtil.Rename(fileinfo.FileName)
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

	fh, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	src, err := fh.Open()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	defer src.Close()
	buf := make([]byte, min(50*1024*1024, fh.Size))

	var fileid string
	if fileinfo.FileId == "" {
		fileid = uuid.MustString()
	} else {
		fileid = fileinfo.FileId
	}

	tempDir := fmt.Sprintf("temp/%v/%v", uid, fileid)
	filePath := fmt.Sprintf("%v/%v", tempDir, fileinfo.ChunkIndex)
	if fileinfo.ChunkIndex == 0 {
		srv.Timer.Add(fileid+contextx.FromUserID(c.Request.Context()), time.Now().Add(time.Minute*30), func() {
			srv.CancelUpload(c, contextx.FromUserID(c.Request.Context()), fileid)
		})
	}

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
			return nil, errors.New("internal error")
		}
	}
	statusMap["fileId"] = fileid
	statusMap["status"] = FILE_STATUS_TRANSFER
	if fileinfo.ChunkIndex == fileinfo.Chunks-1 {
		srv.Timer.Del(fileid + contextx.FromUserID(c.Request.Context()))
		// upload/md5/FileName
		uploadDir := fmt.Sprintf("upload/%v", fileinfo.FileMd5)
		uploadFile := fmt.Sprintf("%v/%v", uploadDir, fileinfo.FileMd5+"."+fileUtil.GetFileExt(fileinfo.FileName))
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
		if file.FileType == define.FileTypeVideo {
			CutFile4Video(fileid, file.FilePath) // 文件切片
			dest := fmt.Sprintf("%s/%s", uploadDir, file.FileMd5+".png")
			go func() {
				CreateCover4Video(file.FilePath, 150, dest)
			}()
			file.FileCover = fmt.Sprintf("%v", file.FileMd5+".png")
		} else if file.FileType == define.FileTypeImage {
			//生成缩略图
			dest := fmt.Sprintf("%s/%s", uploadDir, file.FileMd5+".png")
			go func() {
				CreateCover4Video(file.FilePath, 150, dest)
			}()
			file.FileCover = fmt.Sprintf("%v", file.FileMd5+".png")
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
			return statusMap, nil
		}

		go func() {
			client, err := hdfsUtil.NewClient("172.20.0.2:9000")
			if err != nil {
				log.Default().Printf("error  hdfsUtil newclient error :%v \n", err)
				return
			}

			err = client.CopyDirFromLocal(uploadDir, "/"+uploadDir)
			if err != nil {
				log.Default().Printf("error hdfsUtil CopyDirFromLocal %v to %v err: %v\n", uploadDir, "/"+uploadDir, err)
				return
			}
		}()

		go func() {
			ossClient, err := ossUtil.NewClient()
			if err != nil {
				log.Default().Printf("error  ossUtil newclient error: %v\n", err)
				return
			}

			err = ossClient.CopyDirFromLocal(uploadDir, uploadDir)
			if err != nil {
				log.Default().Printf("error ossUtil CopyDirFromLocal %v to %v error %v\n", uploadDir, "/"+uploadDir, err)
				return
			}
		}()

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

func (f *FileSrv) findAllSubFolderFileIdList(ctx context.Context, fileIdList *[]string, userID, fileID string, delflag int8) {
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
		f.findAllSubFolderFileIdList(ctx, fileIdList, userID, v.FileID, delflag)
	}
}

// 删除文件
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

	//文件删除 应该标志
	for _, e := range fileInfoList {
		go func() {
			f.Timer.Add("file_"+e.FileID+e.UserID, time.Now().Add(time.Hour*24*10), func() {
				//  记录数量

				if e.FileType == define.FileTypeFolder {

				} else {
					count, err := f.Repo.CountFileByMd5(ctx, e.FileMd5)
					if err != nil {
						// 记录日志  /// --->>>
						log.Default().Println("[error] ", err)
						return
					}
					if count == 0 {
						// 错误日志    md5数据有无
						log.Default().Println("[error] ", "文件md5 不存在", e.FileMd5)

					}

					// 如果当前文件
					if count == 1 {
						ossUtil.NewClient()
					}
				}

				f.Repo.DelFiles(ctx, e.UserID, []string{e.FileID})
			})
		}()
		f.findAllSubFolderFileIdList(ctx, &delFileList, uid, e.FileID, define.FileFlagInUse)
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

// 创建视频封面
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
	md5 := name[:strings.LastIndex(name, ".")]

	path := fmt.Sprintf("upload/%s/%s", md5, name)
	cr := cache.NewCacheReader(path)
	reader, err := cr.Read()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Cache-Control", "max-age=2500")
	http.ServeContent(w, r, name, time.Time{}, reader)
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
		fileNameNoSuffix := fmt.Sprintf("%v/%v/%v", "upload", file.FileMd5, tmp)
		cr := cache.NewCacheReader(fileNameNoSuffix)
		reader, err := cr.Read()
		if err != nil {
			return make([]byte, 0), err
		}

		buf, err := io.ReadAll(reader)
		if err != nil {
			return make([]byte, 0), err
		}
		return buf, nil
	} else {
		if (file.FileType) == define.FileTypeVideo {
			prefix := fmt.Sprintf("%v%v", "upload/", file.FileMd5)

			fileNameNoSuffix := fmt.Sprintf("%v%v", prefix, "/index.m3u8")

			reader, err := cache.NewCacheReader(fileNameNoSuffix).Read()
			if err != nil {
				return make([]byte, 0), err
			}
			b, err := io.ReadAll(reader)
			if err != nil {
				return make([]byte, 0), err
			}
			return b, nil

		} else {
			filePath := fmt.Sprintf("%v/%v/%v", "upload", file.FileMd5, file.FileName)
			reader, err := cache.NewCacheReader(filePath).Read()
			if err != nil {
				return make([]byte, 0), err
			}
			b, err := io.ReadAll(reader)
			if err != nil {
				return make([]byte, 0), err
			}

			return b, nil
		}
	}
}

func (f *FileSrv) GetFolderInfo(ctx context.Context, path string, uid string) ([]file.File, error) {
	paths := strings.Split(path, "/")
	var item schema.RequestFileListPage
	item.Path = paths
	item.FolderType = 1
	item.DelFlag = define.FileFlagInUse
	res, err := f.Repo.GetFileList(ctx, uid, &item, false)
	if err != nil {
		return nil, err
	}
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

func (f *FileSrv) ChangeFileFolder(ctx context.Context, uid string, fileIds string, filePid string) error {

	if strings.Contains(fileIds, filePid) {
		return errors.New("")
	}
	if filePid != "0" { //判定父文件夹是否存在
		file, err := f.Repo.GetFileInfo(ctx, filePid, uid)
		if err != nil || file == nil || file.Status != define.FileStatusUsing {
			return errors.New("error") // TODO wait fix
		}
	}

	var item schema.RequestFileListPage

	item.FilePid = filePid
	lists, err := f.Repo.GetFileList(ctx, uid, &item, false) // 找新文件夹下的子文件
	if err != nil {
		log.Println("ChangeFileFolder GetFileList error", err)
		return errors.New("error") // TODO wait fix
	}

	fileNameMap := make(map[string]file.File, 0)
	for _, v := range lists {
		fileNameMap[v.FileName] = v
	}

	item = schema.RequestFileListPage{}
	curs := strings.Split(fileIds, ",")
	item.Path = curs
	lists, _ = f.Repo.GetFileList(ctx, uid, &item, false) // 找要被移动的文件

	for _, ele := range lists {
		if _, ok := fileNameMap[ele.FileName]; ok {
			fileNewName := fileUtil.Rename(ele.FileName)
			ele.FileName = fileNewName
		}
		ele.FilePid = filePid
		if err := f.Repo.UpdateFile(ctx, &ele); err != nil {
			return err
		}
	}

	return nil
}

// 校验文件   // 防止文件跳跃访问      //   环境变量          分享的根目录   用户id   当前文件id
func (f *FileSrv) CheckFootFilePid(ctx context.Context, rootFilePid, userId, fileId string) error {
	fmt.Println("-------------------------------------------------")
	fmt.Printf("rootFilePid: %v\n", rootFilePid)

	if len(rootFilePid) == 0 || len(fileId) == 0 { // 文件id 非法参数
		return errors.New("非法参数")
	}

	if rootFilePid == fileId { // 检验
		return nil
	}

	return f.checkFilePid(ctx, rootFilePid, userId, fileId)
}

// 反响时间最大O(n) 复杂度
func (f *FileSrv) checkFilePid(ctx context.Context, rootFilePid, userId, fileId string) error {
	fileInfo, err := f.Repo.GetFileInfo(ctx, fileId, userId) // TODO 待优化
	if err != nil {
		return err
	}

	if fileInfo == nil {
		return errors.New("文件信息不存在")
	}

	if fileInfo.FileID == "0" {
		return errors.New("非法参数")
	}

	if fileInfo.FilePid == rootFilePid {
		return nil
	}
	return f.checkFilePid(ctx, rootFilePid, userId, fileInfo.FilePid)
}

func (f *FileSrv) SaveShare(ctx context.Context,
	shareRootFilePid,
	shareFileIds, myFolderId,
	shareUserId, currentUserId string) error {

	shareFileIdArray := strings.Split(shareFileIds, ",")
	fmt.Printf("shareFileIdArray: %v\n", shareFileIdArray)

	// 获取当前文件 列表
	query := new(schema.RequestFileListPage)
	query.FilePid = myFolderId
	query.DelFlag = define.FileFlagInUse
	currentFileList, err := f.Repo.GetFileList(ctx, currentUserId, query, false)
	if err != nil {
		return err
	}
	fmt.Printf("currentFileList: %v\n", currentFileList)
	currentFileMap := make(map[string]file.File, 0)

	// 建立一个map 以文件名作为映射
	for _, info := range currentFileList {
		currentFileMap[info.FileName] = info
	}
	fmt.Printf("currentFileMap: %v\n", currentFileMap)

	// 获取当前路径下分享文件列表
	query = new(schema.RequestFileListPage)
	query.Path = shareFileIdArray
	query.DelFlag = define.FileFlagInUse
	shareFileList, err := f.Repo.GetFileList(ctx, shareUserId, query, false)
	if err != nil {
		return err
	}

	fmt.Printf("shareFileList: %v\n", shareFileList)
	fileList := make([]file.File, 0)
	currentTime := time.Now().Format("2006-01-02 15:04:05") // 获取当前的时间
	for _, info := range shareFileList {                    // 便利分享的文件列表
		if _, ok := currentFileMap[info.FileName]; ok { //  检查对应的文件名是否存在
			fileNewName := fileUtil.Rename(info.FileName) // 文件重命令
			info.FileName = fileNewName
		}
		f.findAllSubFileListAndChange(ctx, &fileList, info, shareUserId, currentUserId, currentTime, myFolderId) // 递归 copy  分享文件下的子目录
	}

	for _, file := range fileList { // 文件列表
		if err := f.Repo.UploadFile(ctx, &file); err != nil { // 调用接口保存文件信息  // --- 数据一致性问题  --应该支持回滚操纵
			return err
		}

	}
	return nil
}

// 找出当前分享文件下
func (f *FileSrv) findAllSubFileListAndChange(ctx context.Context, copyFileList *[]file.File, fileInfo file.File, sourceUserId, currentUserID, curentTime string, newFilePid string) {

	// 修改文件信息
	sourceFileId := fileInfo.FileID
	fileInfo.CreateTime = curentTime
	fileInfo.LastUpdateTime = curentTime
	fileInfo.FilePid = newFilePid
	fileInfo.FileID = util.MustString()
	fileInfo.UserID = currentUserID

	*copyFileList = append(*copyFileList, fileInfo) // 文件信息 放入数组里
	if fileInfo.FileType == define.FileTypeFolder { // 当前文件是目录  递归
		query := schema.RequestFileListPage{
			FilePid: sourceFileId,
			DelFlag: define.FileFlagInUse,
		}

		list, err := f.Repo.GetFileList(ctx, sourceUserId, &query, false)
		if err != nil || len(list) == 0 {
			return
		}

		for _, file := range list {
			f.findAllSubFileListAndChange(ctx, copyFileList, file, sourceUserId, currentUserID, curentTime, fileInfo.FileID)
		}
	}

}

// 取消文件上传
func (srv *FileSrv) CancelUpload(ctx context.Context, uid string, fileid string) error {
	path := fmt.Sprintf("temp/%v/%v", uid, fileid) //  取消上传机制
	return os.RemoveAll(path)
}

// 文件加入密码箱
func (srv *FileSrv) UpdateFileSecure(ctx context.Context, uid string, fileid string, status bool) error {
	fileids := strings.Split(fileid, ",")
	return srv.Repo.UpdateFileSecure(ctx, uid, fileids, status)
}

func (srv *FileSrv) GetFileListTotalSize(ctx context.Context, uid string, fileid []string) (uint64, error) {
	var sum uint64
	sum = 0
	md5Set := mapset.NewSet()
	fileMd5 := make([]string, 0)

	for _, id := range fileid {
		srv.findAllSubAllFileIdAndMd5List(ctx, &fileMd5, &md5Set, uid, id, define.FileFlagInUse)
	}
	// 初始化一个 []string 类型的切片
	var stringSlice []string
	md5slice := md5Set.ToSlice()
	for _, elem := range md5slice {
		stringSlice = append(stringSlice, elem.(string))
	}

	files, err := srv.Repo.FindFilesByMd5s(ctx, stringSlice)
	if err != nil {
		return 0, err
	}

	md5ToSize := make(map[string]uint64)
	for _, file := range files {
		md5ToSize[file.FileMd5] = file.FileSize
	}
	for _, md5 := range fileMd5 {
		size, ok := md5ToSize[md5]
		if !ok {
			continue
		}
		sum += size
	}

	return sum, nil
}

// 找出所有md5
func (f *FileSrv) findAllSubAllFileIdAndMd5List(ctx context.Context, fileids *[]string, md5Set *mapset.Set, userID, fileID string, delflag int8) {

	query := schema.RequestFileListPage{
		FilePid: fileID,
		DelFlag: delflag,
	}

	fileLists, err := f.Repo.GetFileList(ctx, userID, &query, false)
	if err != nil || fileLists == nil || len(fileLists) == 0 {
		return
	}

	for _, v := range fileLists {
		if v.FileMd5 != "" {
			*fileids = append(*fileids, v.FileMd5)
			(*md5Set).Add(v.FileMd5)
		}
		f.findAllSubAllFileIdAndMd5List(ctx, fileids, md5Set, userID, v.FileID, delflag)
	}
}

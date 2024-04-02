package service

import (
	"context"
	"log"
	"onlineCLoud/internel/app/dao/file"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/share"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/define"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/pkg/contextx"
	fileUtil "onlineCLoud/pkg/util/file"
	"strings"
	"time"
)

type RecycleSrv struct {
	Repo *file.FileRepo
}

func (f *RecycleSrv) LoadListFiles(ctx context.Context, uid string, pageNo, pageSize int64) (*schema.ListResult, error) {

	var res schema.ListResult
	var p schema.RequestFileListPage
	p.PageNo = int(pageNo)
	p.PageSize = int(pageSize)
	p.DelFlag = define.FileFlagInRecycleBin
	p.OrderBy = "recovery_time desc"
	files, err := f.Repo.GetFileList(ctx, uid, &p, true)
	if err != nil {
		return nil, err
	}
	total, err := f.Repo.GetFileListTotal(ctx, uid, &p)
	if err != nil {
		return nil, err
	}

	res.List = files
	res.Parms = &p.PageParams
	res.PageTotal = (total + int64(p.GetPageSize())/2) / int64(p.GetPageSize())
	res.TotalCount = total
	return &res, nil
}
func (f *RecycleSrv) findAllSubFolderFileList(ctx context.Context, fileIdList *[]string, userID, fileID string, delflag int8) {
	*fileIdList = append(*fileIdList, fileID)

	query := schema.RequestFileListPage{
		FilePid: fileID,
		DelFlag: delflag,
	}

	fields, err := f.Repo.GetFileList(ctx, userID, &query, false)
	if err != nil || fields == nil || len(fields) == 0 {
		return
	}

	for _, v := range fields {
		f.findAllSubFolderFileList(ctx, fileIdList, userID, v.FileID, delflag)
	}

}

func (f *RecycleSrv) DelFiles(ctx context.Context, uid string, fileId string, admin ...bool) error {
	fileIds := strings.Split(fileId, ",")
	query := schema.RequestFileListPage{
		Path:    fileIds,
		DelFlag: define.FileFlagInRecycleBin,
	}

	fileInfoList, _ := f.Repo.GetFileList(ctx, uid, &query, false)
	if fileInfoList == nil || len(fileInfoList) == 0 {
		return nil
	}
	delFileList := make([]string, 0)

	for _, e := range fileInfoList {
		f.findAllSubFolderFileList(ctx, &delFileList, uid, e.FileID, define.FileFlagSoftDeleted)
	}

	delFileList = append(fileIds, delFileList...)

	shareSrv := ShareSrv{Repo: &share.ShareRepo{DB: f.Repo.Db}}
	list, err := shareSrv.LoadShareList(ctx, uid, 0, -1)
	if err != nil {
		return err
	}

	delFileMap := make(map[string]bool)
	for _, delFileID := range delFileList {
		delFileMap[delFileID] = true
	}

	if shareIds, ok := list.List.([]share.Share); ok {
		connIds := make([]string, 0)
		for _, shareItem := range shareIds {
			if delFileMap[shareItem.FileId] {
				connIds = append(connIds, shareItem.ShareId)
			}
		}

		shareSrv.CancelShare(ctx, uid, connIds)
	}

	if err := f.Repo.DelFiles(ctx, uid, delFileList); err != nil {
		return err
	}
	f.delfileToUpdateSpace(ctx, uid)

	return nil
}

func (f *RecycleSrv) delfileToUpdateSpace(ctx context.Context, userid string) {
	// 跟新空间状态
	urv := UserSrv{UserRepo: &user.UserRepo{DB: f.Repo.Db, Rd: redisx.NewClient()}}

	var total uint64
	err := f.Repo.GetTotalUseSpace(ctx, userid, &total)
	if err != nil {
		log.Default().Printf("del file update space error:%v", err)
		return
	}
	urv.UserRepo.UpdateSpace(ctx, contextx.FromUserEmail(ctx), total, true)
}

// 恢复文件-- 》 找当前在回收站的  --》 子目录 --》 父目录更新在根目录下 -- 》 子目录修改状态
func (f *RecycleSrv) RecoverFile(ctx context.Context, uid string, fileIds string) error {
	fildIdArray := strings.Split(fileIds, ",")
	query := schema.RequestFileListPage{
		Path:    fildIdArray,
		DelFlag: define.FileFlagInRecycleBin,
	}
	fileInfoList, err := f.Repo.GetFileList(ctx, uid, &query, false)
	if err != nil {
		return err
	}

	if fileInfoList == nil || len(fileInfoList) == 0 {
		return nil
	}

	FileIDList := make([]string, 0)
	for _, fileinfo := range fileInfoList {
		if fileinfo.FolderType == 1 {
			f.findAllSubFolderFileList(ctx, &FileIDList, uid, fileinfo.FileID, define.FileFlagSoftDeleted) // 子目录
		}
	}
	if len(FileIDList) > 0 {
		if err := f.Repo.UpdateFileDelFlag(ctx, uid, FileIDList, nil, define.FileFlagSoftDeleted, define.FileFlagInUse, ""); err != nil {
			return err
		}
	}

	// 跟新文件状态
	if len(fildIdArray) > 0 {
		if err := f.Repo.UpdateFileDelFlag(ctx, uid, nil, fildIdArray, define.FileFlagInRecycleBin, define.FileFlagInUse, ""); err != nil {
			return err
		}
	}

	rootFileMap := make(map[string]file.File, 0)
	query = schema.RequestFileListPage{
		FilePid: "0",
		DelFlag: define.FileFlagInUse,
	}
	fileinfolist, err := f.Repo.GetFileList(ctx, uid, &query, false)
	if err != nil {
		return err
	}
	for _, v := range fileinfolist {
		rootFileMap[v.FileName] = v
	}

	//移动根目录
	for _, file := range fileInfoList {
		if _, ok := rootFileMap[file.FileName]; ok {
			file.FileName = fileUtil.Rename(file.FileName)
		}
		file.FilePid = "0"
		file.LastUpdateTime = time.Now().Format("2006-01-02 15:04:05")
		file.RecoveryTime = " "
		file.DelFlag = define.FileFlagInUse
		f.Repo.UpdateFile(ctx, &file)

	}

	return nil
}

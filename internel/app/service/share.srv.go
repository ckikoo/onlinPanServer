package service

import (
	"context"
	"fmt"
	"onlineCLoud/internel/app/dao/dto"
	"onlineCLoud/internel/app/dao/file"
	"onlineCLoud/internel/app/dao/share"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/define"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/pkg/contextx"
)

type ShareSrv struct {
	Repo *share.ShareRepo
}

func (f *ShareSrv) LoadShareList(ctx context.Context, uid string, pageNo, pageSize int64) (*schema.ListResult, error) {
	var res schema.ListResult
	var p schema.RequestFileListPage
	p.PageNo = int(pageNo)
	p.PageSize = int(pageSize)
	p.OrderBy = "share_time desc"
	files, err := f.Repo.GetShareList(ctx, uid, &p, true)
	if err != nil {
		return nil, err
	}
	total, err := f.Repo.GetShareListTotal(ctx, uid, &p)
	if err != nil {
		return nil, err
	}

	res.List = files
	res.Parms = &p.PageParams
	res.PageTotal = (total + int64(p.GetPageSize())) / int64(p.GetPageSize())
	res.TotalCount = total
	return &res, nil
}

func (f *ShareSrv) ShareFile(ctx context.Context, share share.Share) error {

	err := f.Repo.Insert(ctx, share)
	return err
}

func (f *ShareSrv) CancelShare(ctx context.Context, uid string, sids []string) error {
	return f.Repo.CancelShare(ctx, uid, sids)
}

func (f *ShareSrv) GetShareLoginInfo(ctx context.Context, shareId string) (*schema.ShareInfo, error) {

	info, err := f.Repo.GetShareInfo(ctx, shareId)
	if err != nil {
		return nil, err
	}

	if info.UserId != contextx.FromUserID(ctx) {
		return nil, nil
	}

	userSrv := UserSrv{UserRepo: &user.UserRepo{DB: f.Repo.DB}}
	userinfo, err := userSrv.GetInfoById(ctx, info.UserId)
	if err != nil {
		return nil, err
	}
	fileSrv := FileSrv{Repo: &file.FileRepo{Db: f.Repo.DB}}
	fileinfo, err := fileSrv.GetFileInfo(ctx, info.FileId, info.UserId)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	shareInfoRes := new(schema.ShareInfo)
	shareInfoRes.ShareTime = info.ShareTime
	shareInfoRes.CurrentUser = info.UserId == contextx.FromUserID(ctx)
	shareInfoRes.Avatar = userinfo.Avatar
	shareInfoRes.UserID = info.UserId
	shareInfoRes.FileID = info.FileId
	shareInfoRes.FileName = fileinfo.FileName
	shareInfoRes.NickName = userinfo.NickName
	if info.ValidType == define.FileShareForverDay {
		shareInfoRes.ExpireTime = "永久"
	} else {
		shareInfoRes.ExpireTime = info.ExpireTime
	}

	return shareInfoRes, nil
}
func (f *ShareSrv) GetShareInfo(ctx context.Context, shareId string) (*schema.ShareInfo, error) {
	info, err := f.Repo.GetShareInfo(ctx, shareId)
	if err != nil {
		return nil, err
	}

	userSrv := UserSrv{UserRepo: &user.UserRepo{DB: f.Repo.DB}}
	userinfo, err := userSrv.GetInfoById(ctx, info.UserId)
	if err != nil {
		return nil, err
	}

	fileSrv := FileSrv{Repo: &file.FileRepo{Db: f.Repo.DB}}
	fileinfo, err := fileSrv.GetFileInfo(ctx, info.FileId, info.UserId)
	if err != nil {
		return nil, err
	}

	shareInfoRes := new(schema.ShareInfo)
	shareInfoRes.ShareTime = info.ShareTime
	shareInfoRes.Avatar = userinfo.Avatar
	shareInfoRes.UserID = info.UserId
	shareInfoRes.FileID = info.FileId
	shareInfoRes.FileName = fileinfo.FileName
	shareInfoRes.NickName = userinfo.NickName
	shareInfoRes.CurrentUser = contextx.FromUserID(ctx) == userinfo.UserID
	fmt.Printf("contextx.FromUserID(ctx): %v\n", contextx.FromUserID(ctx))
	fmt.Printf("userinfo: %v\n", userinfo)
	if info.ValidType == define.FileShareForverDay {
		shareInfoRes.ExpireTime = "永久"
	} else {
		shareInfoRes.ExpireTime = info.ExpireTime
	}

	return shareInfoRes, nil
}

func (f *ShareSrv) GetShareList(ctx context.Context, req *schema.RequestShareListPage) (*schema.ListResult, error) {

	shareInfo, err := f.Repo.GetShareInfo(ctx, req.ShareId)
	if err != nil {
		return nil, err
	}

	var reqFileList schema.RequestFileListPage
	reqFileList.PageNo = req.GetCurrentPage()
	reqFileList.PageSize = req.GetPageSize()
	if len(req.FilePid) != 0 && req.FilePid != "0" { //TODO // 后面修复 无效访问
		reqFileList.FilePid = req.FilePid
	} else {
		reqFileList.Path = []string{shareInfo.FileId}
	}

	fileSrv := FileSrv{Repo: &file.FileRepo{Db: f.Repo.DB}}
	res, err := fileSrv.LoadListFiles(ctx, shareInfo.UserId, &reqFileList)
	if err != nil {
		return nil, err
	}

	return res, nil

}
func (f *ShareSrv) GetFolderInfo(ctx context.Context, shareId string, path string) ([]file.File, error) {

	shareInfo, err := f.Repo.GetShareInfo(ctx, shareId)
	if err != nil {
		return nil, err
	}

	fileSrv := FileSrv{Repo: &file.FileRepo{Db: f.Repo.DB}}

	info, err := fileSrv.GetFolderInfo(ctx, path, shareInfo.UserId)
	if err != nil {
		return nil, err
	}

	return info, nil

}
func (f *ShareSrv) CheckShareCode(ctx context.Context, shareId string, code string) (*dto.SessionWebShareDto, error) {

	shareInfo, err := f.Repo.GetShareInfo(ctx, shareId)
	if err != nil {
		return nil, err
	}

	f.UpdateShareShowCount(ctx, shareId)

	session := new(dto.SessionWebShareDto)
	session.FileId = shareInfo.FileId
	session.Expire = shareInfo.ExpireTime
	session.ShareUserId = shareInfo.UserId
	session.ShareId = shareInfo.ShareId

	return session, nil
}
func (f *ShareSrv) UpdateShareShowCount(ctx context.Context, shareId string) error {

	_, err := f.Repo.UpdateShareShowCount(ctx, shareId)

	return err
}

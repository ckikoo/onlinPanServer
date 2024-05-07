package service

import (
	"context"
	"errors"
	"fmt"
	"onlineCLoud/internel/app/dao/download"
	"time"
)

type DownLoadSrv struct {
	Repo *download.DownloadRepo
}

func (srv *DownLoadSrv) CreateDownLoad(ctx context.Context, userId string, fileId string, path string, code string) error {
	data := download.Download{
		UserId:     userId,
		FileId:     fileId,
		Code:       code,
		CreateTime: time.Now().Unix(),
	}

	f, err := srv.Repo.CreateRecord(ctx, data)
	if err != nil {
		return err
	}
	if f {
		download.CreateRecordUseRedis(ctx, data, path)
		return nil
	}

	return errors.New("创建下载失败")
}

func (srv *DownLoadSrv) FindDownloadByCode(ctx context.Context, code string) (string, error) {

	// 先走redis

	data, err := download.FindRecordByCode(ctx, code)
	if err != nil {
		fmt.Println("err", err)
	}
	if data != "" {
		return data, nil
	}

	// 再走mysql

	data, err = srv.Repo.FindRecordByCode(ctx, code)
	if err != nil {
		return "", err
	}

	return data, nil
}

func (srv *DownLoadSrv) GetALlDownLoad(ctx context.Context, code string) ([]download.Download, error) {

	data, err := srv.Repo.GetAllData(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

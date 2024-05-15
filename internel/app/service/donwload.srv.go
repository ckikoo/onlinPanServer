package service

import (
	"context"
	"errors"
	"fmt"
	"onlineCLoud/internel/app/dao/download"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/user"
	logger "onlineCLoud/pkg/log"
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

	f, err := srv.Repo.CreateRecord(ctx, &data)
	if err != nil {
		return err
	}
	fmt.Printf("data: %v\n", data)
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

func (srv *DownLoadSrv) Delete(ctx context.Context, code string) {
	srv.Repo.Delete(ctx, code)
}

func (srv *DownLoadSrv) GetALlDownLoad(ctx context.Context, code string) ([]download.Download, error) {

	data, err := srv.Repo.GetAllData(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (srv *DownLoadSrv) GetDownLoadSpeed(ctx context.Context, id string) (uint64, error) {
	fmt.Printf("id: %v\n", id)
	if len(id) == 0 {
		return 100, nil
	}
	userRepo := user.UserRepo{
		Rd: redisx.NewClient(),
		DB: srv.Repo.Db,
	}

	limit, err := userRepo.GetUserSpeed(ctx, id)
	if err != nil {
		logger.Log("ERROR", err)
	}

	return limit, nil
}

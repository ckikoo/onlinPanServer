package service

import (
	"context"
	"onlineCLoud/internel/app/dao/enc"
	"onlineCLoud/internel/app/dao/user"
	"time"
)

type EncSrv struct {
	UserRepo *user.UserRepo
	Repo     *enc.EncRepo
}

func (srv *EncSrv) InitPassWord(ctx context.Context, email string, password string) error {
	return srv.UserRepo.UpdateEncPassword(ctx, email, password)
}

func (srv *EncSrv) CheckPassword(ctx context.Context, email string, password string) (bool, error) {
	var user user.User
	err := srv.UserRepo.FindOneByName(ctx, email, &user)
	if err != nil {
		return false, err
	}

	return user.EncPassWord == password, nil
}

func (srv *EncSrv) EnFilePass(ctx context.Context, userId string, fileid string) error {

	err := srv.Repo.AddFile(ctx, &enc.Enc{
		UserID:   userId,
		FileID:   fileid,
		JoinTime: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
func (srv *EncSrv) CheckFileEnc(ctx context.Context, userId string, fileid string, password string) (bool, error) {

	err := srv.Repo.AddFile(ctx, &enc.Enc{
		UserID:   userId,
		FileID:   fileid,
		JoinTime: time.Now(),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

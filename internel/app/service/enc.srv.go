package service

import (
	"context"
	"onlineCLoud/internel/app/dao/user"
)

type EncSrv struct {
	UserRepo *user.UserRepo
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

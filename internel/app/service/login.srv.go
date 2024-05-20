package service

import (
	"context"
	"fmt"
	"onlineCLoud/internel/app/config"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/pkg/auth"
	"onlineCLoud/pkg/errors"
	"onlineCLoud/pkg/util/uuid"
	"time"
)

type LoginSrv struct {
	Auth     auth.Auther
	UserRepo *user.UserRepo
}

func (a *LoginSrv) FindOneByName(ctx context.Context, username string) *user.User {
	var item user.User
	_ = a.UserRepo.FindOneByName(ctx, username, &item)
	return &item
}

func (a *LoginSrv) Login(ctx context.Context, username, password string, isadmin bool) (string, error) {

	var item user.User
	err := a.UserRepo.FindOneByName(ctx, username, &item)
	if err != nil {
		return "", err
	} else if item.Password != password {
		return "", errors.New("邮箱或者密码错误")
	} else if item.Status == 0 {
		return "", errors.New("用户已被禁用，请联系客服申请解禁")
	} else if item.Password == "" {
		return "", errors.New("用户不存在")
	} else if isadmin && !item.Admin {
		return "", errors.New("用户非管理员权限")
	}

	item.LastJoinTime = time.Now()
	err = a.UserRepo.Update(ctx, item.UserID, item)

	return item.UserID, err
}

func (a *LoginSrv) Register(ctx context.Context, User user.User) (bool, error) {
	var item user.User
	err := a.UserRepo.FindOneByName(ctx, User.Email, &item) // 找人
	if err != nil {
		return false, errors.New(errors.ErrInternalServer)
	}

	if item.Password == "" {
		User.UserID = uuid.MustString()
		User.LastJoinTime = time.Now()
		User.CreateTime = User.LastJoinTime
		User.TotalSpace = config.C.File.InitSpaceSize
		User.Status = 1
		fmt.Printf("User: %v\n", User)
		err := a.UserRepo.Create(ctx, &User)
		if err != nil {
			return false, errors.New(errors.ErrInternalServer)
		}
		return true, nil
	}
	return false, errors.New("用户已经存在")
}

func (a *LoginSrv) ResetPasswd(ctx context.Context, email string, password string) (bool, error) {
	err := a.UserRepo.UpdatePassword(ctx, email, password)
	if err != nil {
		fmt.Println("err", err)
		return false, errors.New("修改密码失败")
	}
	return true, nil
}

func (a *LoginSrv) GenerateToken(ctx context.Context, userID string) (*map[string]interface{}, error) {
	tokenInfo, err := a.Auth.GenerateToken(ctx, userID)
	if err != nil {
		return nil, err
	}
	item := make(map[string]interface{}, 0)
	item["AccessToken"] = tokenInfo.GetAccessToken()
	item["ExpiresAt"] = tokenInfo.GetExpiresAt()
	item["TokenType"] = tokenInfo.GetTokenType()
	return &item, nil
}

func (a *LoginSrv) DestoryToken(ctx context.Context, tokenString string) error {
	err := a.Auth.DestroyToken(ctx, tokenString)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

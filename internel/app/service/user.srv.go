package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/pkg/errors"
	"onlineCLoud/pkg/util/hash"
	"os"
	"time"
)

type UserSrv struct {
	UserRepo *user.UserRepo
}

func (a *UserSrv) GetInfoById(ctx context.Context, id string) (*user.User, error) {
	var item user.User
	err := a.UserRepo.FindOneById(ctx, id, &item)

	if err != nil {
		return nil, err
	}

	if item.Password == "" {
		return nil, errors.New("用户信息不存在")
	}
	item.Password = ""
	item.Avatar = ""

	return &item, nil
}

func (a *UserSrv) GetInfo(ctx context.Context, email string) (*user.User, error) {
	var item user.User
	err := a.UserRepo.FindOneByName(ctx, email, &item)

	if err != nil {
		return nil, err
	}

	if item.Password == "" {
		return nil, errors.New("用户信息不存在")
	}
	item.Password = ""
	item.Avatar = ""

	return &item, nil
}

func (a *UserSrv) GetUserSpace(ctx context.Context, email string) map[string]interface{} {

	return a.UserRepo.GetUseSpace(ctx, email)

}

func (a *UserSrv) GetUserSpaceById(ctx context.Context, id string) user.UserSpace {

	return a.UserRepo.GetUserSpaceById(ctx, id)

}
func (a *UserSrv) UpdateSpace(ctx context.Context, email string, add uint64) error {
	return a.UserRepo.UpdateSpace(ctx, email, add)
}
func (a *UserSrv) UpdatePassword(ctx context.Context, email string, old, New string) error {
	var user user.User
	err := a.UserRepo.FindOneByName(ctx, email, &user)
	if err != nil {
		return err
	}
	if user.Password == "" || hash.MD5String(old) != user.Password {
		return errors.New("密码错误")
	}

	return a.UserRepo.UpdatePassword(ctx, email, hash.MD5String(New))
}

func (a *UserSrv) UpdateUserAvatar(ctx context.Context, email string, filename string) error {

	err := a.UserRepo.UpdateUserAvatar(ctx, email, filename)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	return nil
}

func (a *UserSrv) GetUserAvatar(w http.ResponseWriter, r *http.Request, uid string) error {
	var item string
	err := a.UserRepo.FindAvatarByName(r.Context(), uid, &item)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	if item == "" {
		item = "img/default.jpg"
	}
	f, err := os.Open(item)
	if err != nil {
		log.Default().Println(err)
		return err
	}
	defer f.Close()

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Paragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	io.Copy(&content, f)

	http.ServeContent(w, r, item, time.Time{}, bytes.NewReader(content.Bytes()))

	return nil
}

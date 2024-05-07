package service

import (
	"context"
	"onlineCLoud/internel/app/dao/file"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/schema"
)

type AdminSrv struct {
	UserRepo *user.UserRepo
	FileRepo *file.FileRepo
}

func (a *AdminSrv) LoadUserList(ctx context.Context, pageNo int, pageSize int, nickNameFuzzy, status string) (*schema.ListResult, error) {

	q := schema.PageParams{PageNo: pageNo, PageSize: pageSize}

	userList, err := a.UserRepo.LoadUserList(ctx, &q, nickNameFuzzy, status)
	if err != nil {
		return nil, err
	}

	for i, v := range userList {
		v.Avatar = ""
		userList[i] = v
	}

	total, err := a.UserRepo.GetUserListTotal(ctx, &q, nickNameFuzzy, status)

	res := new(schema.ListResult)
	res.PageTotal = (total + int64(pageSize)/2) / int64(pageSize)
	res.Parms = &schema.PageParams{
		PageNo:   pageNo,
		PageSize: pageSize,
	}
	res.List = userList

	res.TotalCount = total

	return res, err
}

func (a *AdminSrv) UpdateUserStatus(ctx context.Context, uid string, status int) (*schema.ListResult, error) {

	userList, err := a.UserRepo.LoadUserList(ctx, &q, nickNameFuzzy, status)
	if err != nil {
		return nil, err
	}

	for i, v := range userList {
		v.Avatar = ""
		userList[i] = v
	}

	total, err := a.UserRepo.GetUserListTotal(ctx, &q, nickNameFuzzy, status)

	res := new(schema.ListResult)
	res.PageTotal = (total + int64(pageSize)/2) / int64(pageSize)
	res.Parms = &schema.PageParams{
		PageNo:   pageNo,
		PageSize: pageSize,
	}
	res.List = userList

	res.TotalCount = total

	return res, err
}

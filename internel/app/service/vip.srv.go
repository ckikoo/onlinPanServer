package service

import (
	"context"
	"onlineCLoud/internel/app/dao/vip"
	"onlineCLoud/internel/app/schema"
)

type VipSrv struct {
	VipRepo *vip.VipRepo
}

func (srv *VipSrv) GetInfo(ctx context.Context, uid string) (interface{}, error) {

	info, err := srv.VipRepo.GetVipInfoByUserID(uid)
	if err != nil {
		return nil, err
	}

	return info, err
}

func (srv *VipSrv) GetVipList(ctx context.Context, pageno, pagesize int, username string) (interface{}, error) {

	var res schema.ListResult

	info, err := srv.VipRepo.LoadVipInfoList(ctx, pageno, pagesize, username)
	if err != nil {
		return nil, err
	}
	res.List = info
	total, err := srv.VipRepo.GetVipListTotal(ctx, pageno, pagesize, username)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total
	res.Parms = &schema.PageParams{
		PageNo:   pageno,
		PageSize: pagesize,
	}
	return res, err
}

func (srv *VipSrv) UpdateTime(ctx context.Context, id, time int, uid string) error {

	err := srv.VipRepo.UpdateTime(ctx, id, time, uid)

	return err
}

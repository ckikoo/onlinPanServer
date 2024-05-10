package service

import (
	"context"
	"onlineCLoud/internel/app/dao/vip"
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

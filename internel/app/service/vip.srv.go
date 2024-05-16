package service

import (
	"context"
	"fmt"
	Pack "onlineCLoud/internel/app/dao/package"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/dao/vip"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/pkg/errors"
	"time"
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
func (srv *VipSrv) AddVip(ctx context.Context, email string, _id, _from, _util int) error {
	if _util <= _from {
		return errors.New("时间不合法")
	}

	userRepo := user.UserRepo{DB: srv.VipRepo.VipDB, Rd: redisx.NewClient()}
	var user user.User
	err := userRepo.FindOneByName(ctx, email, &user)
	if err != nil {
		return err
	}

	if user.Password == "" {
		return errors.New("用户不存在")
	}
	if user.Status == 0 {
		return errors.New("用户已经被封")
	}

	packRepo := Pack.PackageRepo{DB: srv.VipRepo.VipDB}

	info, err := packRepo.FindById(ctx, _id)
	if err != nil {
		return errors.New("出错啦")
	}

	if !info.Show {
		return errors.New("该包已经被禁止啦")
	}

	days := (_util - _from) / 60 / 60 / 24

	err = srv.VipRepo.UpgradeExpireTime(user.UserID, uint32(days))
	if err != nil {
		return errors.New("修改失败")
	}
	fmt.Printf("user: %v\n", user)
	infos, err := srv.VipRepo.GetVipInfoByUserIDAndPackId(user.UserID, uint(_id))
	if err != nil {
		return errors.New("失败")
	}
	fmt.Printf("infos: %v\n", infos)
	if infos != nil && infos.VipPackageID != 0 {
		newT := infos.ActiveUntil.Add(time.Duration(_util) - time.Duration(_from)).Unix()
		err := srv.VipRepo.UpdateTime(ctx, _id, int(newT), infos.UserID)
		if err != nil {
			return errors.New("失败")
		}
	} else {
		fromSeconds := _from / 1000
		fromNanoSeconds := int64((_from % 1000)) * 1e6
		utilSeconds := _util / 1000
		utilNanoSeconds := int64(_util%1000) * 1e6

		fmt.Printf("time.Unix(int64(_from), 0): %v\n", time.Unix(0, int64(_from)))
		err := srv.VipRepo.Insert(ctx, &vip.Vip{
			UserID:       user.UserID,
			VipPackageID: uint(_id),
			ActiveFrom:   time.Unix(int64(fromSeconds), fromNanoSeconds),
			ActiveUntil:  time.Unix(int64(utilSeconds), utilNanoSeconds),
		})
		if err != nil {
			return errors.New("失败")
		}
	}

	return nil

}
func (srv *VipSrv) Delete(ctx context.Context, id int) error {

	err := srv.VipRepo.Delete(ctx, id)

	return err
}

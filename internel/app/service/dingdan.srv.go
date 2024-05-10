package service

import (
	"context"
	"errors"
	"fmt"
	"onlineCLoud/internel/app/dao/dingdan"
	Package "onlineCLoud/internel/app/dao/package"
	"onlineCLoud/internel/app/dao/vip"
	"onlineCLoud/internel/app/schema"
	"time"
)

type DingdanService struct {
	VipRepo     *vip.VipRepo
	PageRepo    *Package.PackageRepo
	DingdanRepo *dingdan.DingdanRepo
}

func (f *DingdanService) GetDingdanList(ctx context.Context, pageNo, pageSize int, page bool, uid string) (*schema.ListResult, error) {

	var res schema.ListResult
	list, err := f.DingdanRepo.GetDingdanList(ctx, pageNo, pageSize, page, uid)
	if err != nil {
		return nil, err
	}
	total, err := f.DingdanRepo.GetDingdanListTotal(ctx, uid)
	if err != nil {
		return nil, err
	}
	res.Parms = &schema.PageParams{
		PageNo:   pageNo,
		PageSize: pageSize,
	}
	res.List = list
	res.TotalCount = total
	return &res, nil
}

func (f *DingdanService) FindById(ctx context.Context, id int) (*dingdan.Dingdan, error) {
	dingdan, err := f.DingdanRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return dingdan, nil
}

func (f *DingdanService) Buy(ctx context.Context, uid string, id int) error {
	pinfo, err := f.PageRepo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if pinfo.Show == false {
		return errors.New("该套餐已下架")
	}

	fmt.Printf("pinfo: %+v\n", pinfo)

	f.VipRepo.UpgradeExpireTime(uid, pinfo.ExpireDays)
	err = f.DingdanRepo.Insert(ctx, dingdan.Dingdan{UserId: uid, PackageId: id})
	if err != nil {
		return err
	}

	ex, err := f.VipRepo.CheckExists(uid, uint(id))
	if err != nil {
		return err
	}
	if ex {
		return nil
	}

	err = f.VipRepo.Insert(ctx, &vip.Vip{
		UserID:       uid,
		VipPackageID: uint(id),
		ActiveFrom:   time.Now(),
		ActiveUntil:  time.Now().Add(time.Hour * 24 * time.Duration(pinfo.ExpireDays)),
	})
	return err
}

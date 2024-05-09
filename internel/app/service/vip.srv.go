package service

import (
	"context"
	"onlineCLoud/internel/app/dao/vip"
	"onlineCLoud/internel/app/schema"
)

type VipService struct {
	Repo *vip.VipRepo
}

func (f *VipService) GetVipList(ctx context.Context) (*schema.ListResult, error) {

	var res schema.ListResult
	list, err := f.Repo.GetVipList(ctx, false)
	if err != nil {
		return nil, err
	}
	total, err := f.Repo.GetVipListTotal(ctx)
	if err != nil {
		return nil, err
	}

	res.List = list
	res.TotalCount = total
	return &res, nil
}

func (f *VipService) DelVips(ctx context.Context, ids int) error {

	err := f.Repo.Delete(ctx, ids)
	if err != nil {
		return err
	}

	return nil
}

func (f *VipService) AddVip(ctx context.Context, vip vip.Vip) error {

	err := f.Repo.Insert(ctx, vip)
	return err

}

func (f *VipService) Update(ctx context.Context, id int, vip vip.Vip) error {
	err := f.Repo.Update(ctx, id, vip)
	if err != nil {
		return err
	}

	return nil

}
func (f *VipService) UpdateStatus(ctx context.Context, id int, status bool) error {
	err := f.Repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	return nil

}

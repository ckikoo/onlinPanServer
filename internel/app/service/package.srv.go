package service

import (
	"context"
	"onlineCLoud/internel/app/dao/dingdan"
	Package "onlineCLoud/internel/app/dao/package"
	"onlineCLoud/internel/app/schema"
)

type PageService struct {
	Repo        *Package.PackageRepo
	DingDanRepo *dingdan.DingdanRepo
}

func (f *PageService) GetPackageList(ctx context.Context, status bool) (*schema.ListResult, error) {

	var res schema.ListResult
	list, err := f.Repo.GetPageList(ctx, status, false)
	if err != nil {
		return nil, err
	}
	total, err := f.Repo.GetPageListTotal(ctx, status)
	if err != nil {
		return nil, err
	}

	res.List = list
	res.TotalCount = total
	return &res, nil
}

func (f *PageService) DelPackages(ctx context.Context, ids int) error {

	err := f.Repo.Delete(ctx, ids)
	if err != nil {
		return err
	}

	return nil
}

func (f *PageService) AddPackage(ctx context.Context, Package Package.Package) error {

	err := f.Repo.Insert(ctx, Package)
	return err

}

func (f *PageService) Update(ctx context.Context, id int, Package Package.Package) error {
	err := f.Repo.Update(ctx, id, Package)
	if err != nil {
		return err
	}

	return nil

}
func (f *PageService) UpdateStatus(ctx context.Context, id int, status bool) error {
	err := f.Repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	return nil

}
func (f *PageService) FindById(ctx context.Context, id int) (*Package.Package, error) {
	temp, err := f.Repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return temp, nil
}

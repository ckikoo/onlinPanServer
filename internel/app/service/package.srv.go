package service

import (
	"context"
	"onlineCLoud/internel/app/dao/pkg"
	"onlineCLoud/pkg/util/uuid"
	"time"
)

type PackageService struct {
	Repo *pkg.PkgRepo
}

func (srv *PackageService) GetPackInfo(ctx context.Context) (any, error) {

	res, err := srv.Repo.GetPackInfo(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (srv *PackageService) CheckExists(ctx context.Context, packId string) (any, error) {
	info, err := srv.Repo.GetPackInfoByID(ctx, packId)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (srv *PackageService) BuySpace(ctx context.Context, uid string, info pkg.Pkg) (bool, error) {
	size := 1024 * 1024 * 1024 * uint64(info.Size)
	pkg := pkg.BuySpace{
		SpaceId:   uuid.MustString(),
		UserId:    uid,
		Size:      size,
		CreatedAt: time.Now(),
		UntilAt:   time.Now().Add(time.Hour * 24 * 30),
	}
	res, err := srv.Repo.BuySpace(ctx, uid, pkg)
	if err != nil || res == false {
		return res, err
	}

	return true, nil

}

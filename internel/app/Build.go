package app

import (
	"context"
	"onlineCLoud/internel/app/api"
	"onlineCLoud/internel/app/dao/file"
	"onlineCLoud/internel/app/dao/mailx"
	"onlineCLoud/internel/app/dao/pkg"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/share"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/router"
	"onlineCLoud/internel/app/service"
)

func BuildInjector() (*Injector, func(), error) {
	auther, cleanup, err := InitAuth()
	if err != nil {
		return nil, nil, err
	}

	db, cleanup2, err := InitGormDB()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	cleanup3 := mailx.Init()

	UserRepo := user.UserRepo{
		DB: db,
		Rd: redisx.NewClient(),
	}
	loginSrv := service.LoginSrv{
		Auth:     auther,
		UserRepo: &UserRepo,
	}
	loginApi := api.LoginAPI{
		LoginSrv: &loginSrv,
	}

	userSrv := service.UserSrv{
		UserRepo: &UserRepo,
	}
	userApi := api.UserAPI{
		UserSrv: &userSrv,
	}

	fileRepo := file.FileRepo{
		Db: db,
	}
	FileSrv := service.FileSrv{
		Repo: &fileRepo,
	}
	fileApi := api.FileApi{
		FileSrv: &FileSrv,
	}

	RecycleSrv := service.RecycleSrv{
		Repo: &fileRepo,
	}
	recycleApi := api.RecycleApi{
		RecycleSrv: &RecycleSrv,
	}

	AdminSrv := service.AdminSrv{
		UserRepo: &UserRepo,
		FileRepo: &fileRepo,
	}
	AdminApi := api.AdminApi{
		AdminSrv: &AdminSrv,
	}
	ShareSrv := service.ShareSrv{
		Repo: &share.ShareRepo{DB: db},
	}
	ShareApi := api.ShareApi{
		ShareSrv: &ShareSrv,
	}

	packageApi := api.PackageApi{
		Srv: &service.PackageService{Repo: &pkg.PkgRepo{DB: db, RD: redisx.NewClient()}},
	}

	WebShareApi := api.WebShareApi{
		ShareSrv: &ShareSrv,
		FileSrv:  &FileSrv,
	}

	go file.Init(context.Background(), db, redisx.NewClient())
	routerRouter := &router.Router{
		Auth:        auther,
		LoginAPI:    &loginApi,
		UserApi:     &userApi,
		FileApi:     &fileApi,
		RecycleApi:  &recycleApi,
		ShareApi:    &ShareApi,
		AdminApi:    &AdminApi,
		PackageApi:  &packageApi,
		WebShareApi: &WebShareApi,
	}

	engine := InitGinEngine(routerRouter)

	injector := &Injector{
		Engine: engine,
		Auth:   auther,
	}

	return injector, func() {
		cleanup()
		cleanup2()
		cleanup3()
	}, nil

}

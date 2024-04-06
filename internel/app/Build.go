package app

import (
	"context"
	"onlineCLoud/internel/app/api"
	"onlineCLoud/internel/app/dao/file"
	"onlineCLoud/internel/app/dao/mailx"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/share"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/define"
	"onlineCLoud/internel/app/router"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/timer"
	"time"
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
	timer, timerClean := timer.NewTimerManager()

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
		Repo:  &fileRepo,
		Timer: timer,
	}
	fileApi := api.FileApi{
		FileSrv: &FileSrv,
	}
	EncSrv := service.EncSrv{
		UserRepo: &UserRepo,
	}
	EncApi := api.EncAPI{
		EncSrv:  &EncSrv,
		FileSrv: &FileSrv,
	}

	RecycleSrv := service.RecycleSrv{
		Repo:  &fileRepo,
		Timer: timer,
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

	WebShareApi := api.WebShareApi{
		ShareSrv: &ShareSrv,
		FileSrv:  &FileSrv,
	}

	routerRouter := &router.Router{
		Auth:        auther,
		LoginAPI:    &loginApi,
		UserApi:     &userApi,
		FileApi:     &fileApi,
		RecycleApi:  &recycleApi,
		ShareApi:    &ShareApi,
		AdminApi:    &AdminApi,
		WebShareApi: &WebShareApi,
		EncAPI:      &EncApi,
	}

	go func() {
		list, _ := fileRepo.GetFileList(context.Background(), "*", &schema.RequestFileListPage{DelFlag: define.FileFlagInRecycleBin}, false)
		for _, file := range list {
			joinTime, _ := time.Parse("2006-01-02 15:04:05", file.RecoveryTime)
			EndTime := joinTime.Add(time.Hour * 24 * 10)
			timer.Add("file_"+file.FileID+file.UserID, EndTime, func() {
				RecycleSrv.DelFiles(context.Background(), file.UserID, file.FileID)
			})
		}
	}()

	engine := InitGinEngine(routerRouter)

	injector := &Injector{
		Engine: engine,
		Auth:   auther,
	}

	return injector, func() {
		cleanup()
		cleanup2()
		cleanup3()
		timerClean()
	}, nil

}

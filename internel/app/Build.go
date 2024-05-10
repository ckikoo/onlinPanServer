package app

import (
	"context"
	"onlineCLoud/internel/app/api"
	"onlineCLoud/internel/app/api/admin"
	"onlineCLoud/internel/app/dao/dingdan"
	"onlineCLoud/internel/app/dao/download"
	"onlineCLoud/internel/app/dao/enc"
	"onlineCLoud/internel/app/dao/file"
	workOrder "onlineCLoud/internel/app/dao/gongdan"
	"onlineCLoud/internel/app/dao/mailx"
	Package "onlineCLoud/internel/app/dao/package"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/share"
	"onlineCLoud/internel/app/dao/user"
	"onlineCLoud/internel/app/dao/vip"
	"onlineCLoud/internel/app/define"
	"onlineCLoud/internel/app/router"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/internel/app/service"
	logger "onlineCLoud/pkg/log"
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
	timer := timer.GetInstance()

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
	downloadRepo := download.DownloadRepo{
		Db: db,
	}

	downloadSrv := service.DownLoadSrv{
		Repo: &downloadRepo,
	}

	fileApi := api.FileApi{
		FileSrv:     &FileSrv,
		DownLoadSrv: &downloadSrv,
	}
	EncSrv := service.EncSrv{
		Repo: &enc.EncRepo{
			Db: db,
		},
		UserRepo: &UserRepo,
	}
	EncApi := api.EncAPI{
		EncSrv:  &EncSrv,
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
	AdminApi := admin.AdminApi{
		AdminSrv: &AdminSrv,
	}
	AdminLoginAPI := admin.AdminLoginAPI{
		LoginSrv: &loginSrv,
	}
	ShareSrv := service.ShareSrv{
		Repo: &share.ShareRepo{DB: db},
	}
	ShareApi := api.ShareApi{
		ShareSrv: &ShareSrv,
	}

	WebShareApi := api.WebShareApi{
		UserSrv:  &userSrv,
		ShareSrv: &ShareSrv,
		FileSrv:  &FileSrv,
	}

	work := api.WorkOrderApi{
		Srv: service.WorkOrderSrv{
			Repo: &workOrder.WorkOrderRepo{
				DB: db,
			},
		},
	}

	AdminOrder := admin.AdminOrderApi{
		Srv: &service.WorkOrderSrv{
			Repo: &workOrder.WorkOrderRepo{
				DB: db,
			},
		},
	}

	pageSrv := service.PageService{
		Repo: &Package.PackageRepo{
			DB: db,
		},
		DingDanRepo: &dingdan.DingdanRepo{
			DB: db,
		},
	}

	adminPackage := admin.PackageApi{Srv: &pageSrv}
	pa := api.PageApi{
		Srv: &pageSrv,
	}

	download := api.DownLoadApi{
		Srv:     &downloadSrv,
		FileSrv: &FileSrv,
	}
	routerRouter := &router.Router{
		Auth:          auther,
		LoginAPI:      &loginApi,
		UserApi:       &userApi,
		FileApi:       &fileApi,
		RecycleApi:    &recycleApi,
		ShareApi:      &ShareApi,
		AdminLoginApi: &AdminLoginAPI,
		AdminApi:      &AdminApi,
		WebShareApi:   &WebShareApi,
		EncAPI:        &EncApi,
		WorkOrder:     &work,
		AdminOrder:    &AdminOrder,
		DownLoad:      &download,
		AdminPackage:  &adminPackage,
		PageApi:       &pa,
		Dingdan: &api.DingdanApi{
			Srv: &service.DingdanService{
				PageRepo: pageSrv.Repo,
				DingdanRepo: &dingdan.DingdanRepo{
					DB: db,
				},
				VipRepo: &vip.VipRepo{
					VipDB: db,
				},
			},
		},
		Vip: &api.VipAPI{
			VipSrv: &service.VipSrv{
				VipRepo: &vip.VipRepo{
					VipDB: db,
				},
			},
		},
		AdminVip: &admin.VipAPI{
			VipSrv: &service.VipSrv{
				VipRepo: &vip.VipRepo{
					VipDB: db,
				},
			},
		},
	}

	// 对回收站定时删除
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

	go InitTimer()
	engine := InitGinEngine(routerRouter)

	injector := &Injector{
		Engine: engine,
		Auth:   auther,
	}

	// 初始化日志
	logger.InitLogger()

	return injector, func() {
		cleanup()
		cleanup2()
		cleanup3()
		timer.Close()
		logger.Close()
	}, nil

}

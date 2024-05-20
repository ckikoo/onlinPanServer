package dao

import (
	"onlineCLoud/internel/app/config"
	"onlineCLoud/internel/app/dao/dingdan"
	"onlineCLoud/internel/app/dao/download"
	"onlineCLoud/internel/app/dao/enc"
	"onlineCLoud/internel/app/dao/file"
	workOrder "onlineCLoud/internel/app/dao/gongdan"
	Package "onlineCLoud/internel/app/dao/package"
	"onlineCLoud/internel/app/dao/recycle"
	"onlineCLoud/internel/app/dao/share"
	"onlineCLoud/internel/app/dao/vip"
	"os/user"
	"strings"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	// TODO 数据库创建表
	err := db.AutoMigrate(
		new(user.User),
		new(file.File),
		new(share.Share),
		new(workOrder.WorkOrder),
		new(download.Download),
		new(enc.Enc),
		new(Package.Package),
		new(dingdan.Dingdan),
		new(vip.Vip),
		new(recycle.Recycle),
	)

	return err

}

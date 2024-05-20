package recycle

import (
	"context"
	"onlineCLoud/internel/app/dao/util"

	"gorm.io/gorm"
)

type Recycle struct {
	FileID       string `json:"fileId"     form:"fileId" gorm:"column:file_id;type:varchar(36) ;primaryKey;"` //文件编号
	UserID       string `json:"userId"     form:"userId" gorm:"column:user_id;type:varchar(36);primaryKey"`   //用户编号
	RecoveryTime string `json:"recoveryTime" form:"recoveryTime" gorm:"column:recovery_time"`                 //进入回收站时间
}

func GetRecycleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Recycle))
}

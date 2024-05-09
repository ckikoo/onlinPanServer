package enc

import (
	"context"
	"onlineCLoud/internel/app/dao/util"
	"time"

	"gorm.io/gorm"
)

type Enc struct {
	FileID   string    `json:"fileId"    form:"fileId" gorm:"column:file_id;type:varchar(36);primaryKey"` // 文件编号
	UserID   string    `json:"userId"    form:"userId" gorm:"column:user_id;type:varchar(36);primaryKey"` // 用户编号
	JoinTime time.Time `json:"joinTime"  gorm:"column:join_time;type:datetime"`                           // 加入时间
}

func GetEncDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Enc))
}

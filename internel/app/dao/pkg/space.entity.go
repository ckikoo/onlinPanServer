package pkg

import (
	"context"
	"onlineCLoud/internel/app/dao/util"
	"time"

	"gorm.io/gorm"
)

type BuySpace struct {
	SpaceId   string `json:"spaceId" form:"spaceId" gorm:"column:space_id;type:varchar(36);primaryKey;"`
	UserId    string `json:"userId" form:"userId" gorm:"column:user_id;type:varchar(36);"`
	Size      uint64 `json:"size" form:"size" gorm:"column:size;type:bigint(20) "`
	CreatedAt time.Time
	UntilAt   time.Time
}

func GetSpaceDB(ctx context.Context, old *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, old, new(BuySpace))
}

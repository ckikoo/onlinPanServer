package Pack

import (
	"context"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	SpeedLimit uint64  `gorm:"column:speedLimit" json:"speedLimit" form:"speedLimit"` // 下载限制速度
	SpaceSize  uint64  `gorm:"column:spaceSize" json:"spaceSize" form:"spaceSize"`    // 容量限制
	Show       bool    `gorm:"column:show" json:"show" form:"show'"`                  // 状态，可能表示是否激活
	PageName   string  `gorm:"column:pageName" json:"pageName" form:"pageName"`       // 页面名称
	Price      float32 `gorm:"column:price" json:"price" form:"price"`                // 价格
	ExpireDays uint32  `gorm:"column:days" json:"days" form:"days"`                   // 有效时间
}

func GETPackageDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.WithContext(ctx).Model(&Package{})
}

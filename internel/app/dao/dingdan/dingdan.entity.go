package dingdan

import (
	"context"

	"gorm.io/gorm"
)

type Dingdan struct {
	gorm.Model
	UserId    string `gorm:"column:user_id;index" json:"user_id"`
	PackageId int    `gorm:"column:package_id;index" json:"package_id"`
}

func GetDingdanDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.WithContext(ctx).Model(&Dingdan{})
}

package vip

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Vip struct {
	gorm.Model
	UserID       string
	VipPackageID uint
	ActiveFrom   time.Time
	ActiveUntil  time.Time
}

func GetVipDb(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.WithContext(ctx).Model(&Vip{})
}

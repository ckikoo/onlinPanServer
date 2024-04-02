package util

import (
	"context"
	"encoding/json"
	"onlineCLoud/internel/app/schema"

	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint64    `json:"id" gorm:"primaryKey;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ParseJson(ctx context.Context, str string, out interface{}) error {
	return json.Unmarshal([]byte(str), out)
}
func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	//TODO 悲观锁

	return defDB
}

func GetDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return GetDB(ctx, defDB).Model(m)
}

func FindOne(ctx context.Context, db *gorm.DB, out interface{}) (bool, error) {
	result := db.First(out)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func WrapPageQuery(ctx context.Context, db *gorm.DB, pp *schema.PageParams, items interface{}, page bool) error {
	if page {
		current, pageSize := pp.GetCurrentPage(), pp.GetPageSize()
		if current > 0 && pageSize > 0 {
			db = db.Offset((current - 1) * pageSize).Limit(pageSize)
		} else if pageSize > 0 {
			db = db.Limit(pageSize)
		}
	}

	err := db.Find(items).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return nil
}

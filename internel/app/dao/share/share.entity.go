package share

import (
	"context"

	"gorm.io/gorm"
)

type Share struct {
	ShareId    string `json:"shareId" form:"shareId" gorm:"column:share_id;type:varchar(36) ;primaryKey"`
	FileId     string `json:"fileId" form:"fileId" gorm:"column:file_id;type:varchar(36) ;"`
	UserId     string `json:"userId" form:"userId" gorm:"column:user_id;type:varchar(36) ;"`
	ExpireTime string `json:"expireTime"`
	ShareTime  string `json:"shareTime"`
	Code       string `json:"code" form:"code" gorm:"column:code;type:varchar(5) ;"`
	ValidType  int8   `json:"validType" form:"validType" gorm:"column:valid_type;type:tinyint(1);"`
	ShowCount  uint32 `json:"showCount" form:"showCount" gorm:"column:show_count;type:int(11);"`
}

func GetShareDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.Model(&Share{})
}

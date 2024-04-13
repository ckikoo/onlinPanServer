package workOrder

import (
	"context"

	"gorm.io/gorm"
)

type WorkOrder struct {
	WorkOrderId string `json:"workOrderId" form:"workOrderId" gorm:"column:workOrderId;type:varchar(36);primaryKey"`
	Title       string `json:"title" form:"title" gorm:"column:title"`
	Content     string `json:"content" form:"content" gorm:"column:content"`
	UserId      string `json:"userId" form:"userId" gorm:"column:user_id;index"`
	Status      bool   `json:"status" form:"status" gorm:"column:status"`
	ReplyCotent string `json:"replyCotent" form:"replyCotent" gorm:"column:reply_content"`
	AdminId     string `json:"adminId" form:"adminId" gorm:"column:admin_id"`
	CreateTime  uint64 `json:"create_time" form:"create_time" gorm:"create_time"`
	DoneTime    uint64 `json:"done_time" form:"done_time" gorm:"done_time"`
}

func GetWorkOrderDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.Model(&WorkOrder{})
}

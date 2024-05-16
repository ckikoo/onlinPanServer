package workOrder

import (
	"context"
	"log"
	"onlineCLoud/internel/app/dao/util"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/pkg/errors"

	"gorm.io/gorm"
)

type WorkOrderRepo struct {
	DB *gorm.DB
}
type order struct {
	WorkOrderId string `json:"workOrderId" form:"workOrderId" gorm:"column:workOrderId;type:varchar(36);primaryKey"`
	Title       string `json:"title" form:"title" gorm:"column:title"`
	Content     string `json:"content" form:"content" gorm:"column:content"`
	UserId      string `json:"userId" form:"userId" gorm:"column:user_id;index"`
	Email       string `json:"email" form:"email" gorm:"column:email"`
	Status      bool   `json:"status" form:"status" gorm:"column:status"`
	ReplyCotent string `json:"replyCotent" form:"replyCotent" gorm:"column:reply_content"`
	AdminId     string `json:"adminId" form:"adminId" gorm:"column:admin_id"`
	CreateTime  uint64 `json:"create_time" form:"create_time" gorm:"create_time"`
	DoneTime    uint64 `json:"done_time" form:"done_time" gorm:"done_time"`
}

func (a *WorkOrderRepo) LoadWordOrderList(ctx context.Context, p *schema.PageParams, userId string, status string) ([]order, error) {
	var list []order
	db := GetWorkOrderDB(ctx, a.DB).
		Joins("join tb_user on tb_user.user_id = tb_work_order.user_id").
		Select("tb_user.email, tb_work_order.*")

	// 添加查询条件
	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}
	if status != "*" {
		db = db.Where("tb_work_order.status = ?", status)
	}
	db = db.Order("tb_work_order.create_time desc")
	db = db.Order("tb_work_order.status asc")
	// 分页查询
	err := util.WrapPageQuery(ctx, db, p, &list, true)
	return list, err
}

func (a *WorkOrderRepo) GetWordListTotal(ctx context.Context, p *schema.PageParams, userID string, status string) (int64, error) {
	db := GetWorkOrderDB(ctx, a.DB)

	// 添加查询条件
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	if status != "*" {
		db = db.Where("tb_work_order.status = ?", status)
	}

	// 计算总数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		log.Println(err)
		return 0, err // 返回错误
	}
	return total, nil
}

func (a *WorkOrderRepo) Create(ctx context.Context, item *WorkOrder) error {
	// 执行创建操作
	result := GetWorkOrderDB(ctx, a.DB).Create(item)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to create work order") // 返回错误
	}
	return nil
}
func (a *WorkOrderRepo) FindWorkOrderById(ctx context.Context, id string) (*WorkOrder, error) {
	item := new(WorkOrder)
	// 执行创建操作
	result := GetWorkOrderDB(ctx, a.DB).Where(&WorkOrder{WorkOrderId: id}).First(item)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(result.Error, "failed to create work order") // 返回错误
	}
	return item, nil
}

func (a *WorkOrderRepo) Update(ctx context.Context, id string, item *WorkOrder) error {
	// 执行更新操作
	result := GetWorkOrderDB(ctx, a.DB).Where("workOrderId=?", id).Updates(item)
	if result.Error != nil {
		return errors.Wrap(result.Error, "workOrder to update work order") // 返回错误
	}
	return nil
}

func (a *WorkOrderRepo) Delete(ctx context.Context, id string, orderId string) error {
	// 执行删除操作
	db := GetWorkOrderDB(ctx, a.DB)

	// 执行删除操作
	result := db.Where(&WorkOrder{UserId: id, WorkOrderId: orderId}).Delete(&WorkOrder{})
	if result.Error != nil {
		// 错误处理
		return errors.Wrap(result.Error, "删除错误")
	}

	// 验证删除是否成功
	if result.RowsAffected == 0 {
		return errors.New("no work order deleted")
	}

	return nil
}

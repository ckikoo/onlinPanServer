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

func (a *WorkOrderRepo) LoadWordOrderList(ctx context.Context, p *schema.PageParams, userId string, status string) ([]WorkOrder, error) {
	var list []WorkOrder
	db := GetWorkOrderDB(ctx, a.DB)

	// 添加查询条件
	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}
	if status != "*" {
		db = db.Where("status = ?", status)
	}
	db = db.Order("create_time desc")
	db = db.Order("status asc")
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
		db = db.Where("status = ?", status)
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

func (a *WorkOrderRepo) Delete(ctx context.Context, id string, item WorkOrder) error {
	// 执行删除操作
	result := GetWorkOrderDB(ctx, a.DB).Where("WorkOrderId=?", id).Delete(item)
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete work order") // 返回错误
	}
	return nil
}

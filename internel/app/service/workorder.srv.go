package service

import (
	"context"
	"errors"
	workOrder "onlineCLoud/internel/app/dao/gongdan"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/pkg/util/uuid"
	"time"
)

type WorkOrderSrv struct {
	Repo *workOrder.WorkOrderRepo
}

func (srv *WorkOrderSrv) LoadWorkListList(ctx context.Context, PageNo int, PageSize int, userId string, status string) (*schema.ListResult, error) {
	q := schema.PageParams{PageNo: PageNo, PageSize: PageSize}

	wrks, err := srv.Repo.LoadWordOrderList(ctx, &q, userId, status)
	if err != nil {
		return nil, err
	}
	total, err := srv.Repo.GetWordListTotal(ctx, &q, userId, status)
	if err != nil {
		return nil, err
	}

	res := new(schema.ListResult)
	res.PageTotal = (total + int64(PageNo)/2) / int64(PageSize)
	res.Parms = &schema.PageParams{
		PageNo:   PageNo,
		PageSize: PageSize,
	}
	res.List = wrks

	res.TotalCount = total

	return res, nil
}

func (srv *WorkOrderSrv) CreateWorkOrder(ctx context.Context, uid string, title, content string) error {
	item := workOrder.WorkOrder{
		WorkOrderId: uuid.MustString(),
		UserId:      uid,
		Title:       title,
		Content:     content,
		CreateTime:  uint64(time.Now().Unix()),
		Status:      false,
	}

	return srv.Repo.Create(ctx, &item)

}

func (srv *WorkOrderSrv) UpdateWorkOrder(ctx context.Context, uid string, workId string, title, content string) error {
	item, err := srv.Repo.FindWorkOrderById(ctx, workId)
	if err != nil {
		return err
	}
	if len(item.UserId) == 0 || (uid != "*" && item.UserId != uid) {
		return errors.New("工单不存在")
	}
	if item.Status {
		return errors.New("工单已经处理")
	}

	item.Title = title
	item.Content = content

	return srv.Repo.Update(ctx, item.WorkOrderId, item)
}

func (srv *WorkOrderSrv) AdminUpdateWorkOrder(ctx context.Context, uid string, workId string, title, content string, reply string) error {
	item, err := srv.Repo.FindWorkOrderById(ctx, workId)
	if err != nil {
		return err
	}
	item.DoneTime = uint64(time.Now().Unix())
	item.AdminId = uid
	item.Status = true
	item.Title = title
	item.Content = content
	item.ReplyCotent = reply

	return srv.Repo.Update(ctx, item.WorkOrderId, item)
}
func (srv *WorkOrderSrv) DeleOrederSrv(ctx context.Context, uid string, workId string) error {
	err := srv.Repo.Delete(ctx, uid, workId)
	if err != nil {
		return err
	}
	return nil
}

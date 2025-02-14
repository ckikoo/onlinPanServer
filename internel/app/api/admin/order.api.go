package admin

import (
	"fmt"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminOrderApi struct {
	Srv *service.WorkOrderSrv
}

func (a *AdminOrderApi) LoadWorkList(c *gin.Context) {
	ctx := c.Request.Context()
	pageNo := c.Request.PostFormValue("pageNo")
	pageSize := c.Request.PostFormValue("pageSize")
	status := c.PostForm("status")
	fmt.Printf("status: %v\n", status)
	if status == "" {
		status = "*"
	}
	if pageNo == "" {
		pageNo = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	PageNo, err := strconv.ParseInt(pageNo, 10, 64)
	if err != nil {
		ginx.ResFail(c)
		return
	}
	PageSize, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		ginx.ResFail(c)
		return
	}
	fmt.Printf("status: %v\n", status)
	if status == "1" {

		res, err := a.Srv.LoadWorkListList(ctx, int(PageNo), int(PageSize), "", "true")
		if err != nil {
			ginx.ResFail(c)
			return
		}
		ginx.ResOkWithData(c, res)
	} else if status == "0" {
		res, err := a.Srv.LoadWorkListList(ctx, int(PageNo), int(PageSize), "", "false")
		if err != nil {
			ginx.ResFail(c)
			return
		}
		ginx.ResOkWithData(c, res)
	} else {
		res, err := a.Srv.LoadWorkListList(ctx, int(PageNo), int(PageSize), "", status)
		if err != nil {
			ginx.ResFail(c)
			return
		}
		ginx.ResOkWithData(c, res)
	}

}

func (a *AdminOrderApi) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Request.PostFormValue("id")
	content := c.Request.PostFormValue("content")
	title := c.Request.PostFormValue("title")
	reply := c.Request.PostFormValue("reply")
	if len(id) == 0 || len(content) == 0 || len(title) == 0 || len(reply) == 0 {
		ginx.ResFail(c)
		return
	}

	err := a.Srv.AdminUpdateWorkOrder(ctx, contextx.FromUserID(ctx), id, title, content, reply)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithMessage(c, "修改成功")
}
func (a *AdminOrderApi) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Request.PostFormValue("userId")
	workOrderId := c.Request.PostFormValue("workOrderId")

	if len(id) == 0 || len(workOrderId) == 0 {
		ginx.ResFail(c)
		return
	}

	err := a.Srv.DeleOrederSrv(ctx, id, workOrderId)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithMessage(c, "修改成功")
}

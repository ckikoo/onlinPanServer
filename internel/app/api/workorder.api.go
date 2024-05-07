package api

import (
	"fmt"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkOrderApi struct {
	Srv service.WorkOrderSrv
}

func (a *WorkOrderApi) LoadWorkList(c *gin.Context) {
	ctx := c.Request.Context()
	pageNo := c.Request.PostFormValue("pageNo")
	pageSize := c.Request.PostFormValue("pageSize")
	status := c.DefaultPostForm("status", "*")
	if pageNo == "" && pageSize == "" {
		pageNo = "1"
		pageSize = "20"
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

	res, err := a.Srv.LoadWorkListList(ctx, int(PageNo), int(PageSize), contextx.FromUserID(ctx), status)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, res)
}

func (a *WorkOrderApi) UpdateWorkOrder(c *gin.Context) {
	ctx := c.Request.Context()
	content := c.Request.PostFormValue("content")
	title := c.Request.PostFormValue("title")
	workid := c.Request.PostFormValue("id")

	if len(content) == 0 || len(title) == 0 || len(workid) == 0 {
		ginx.ResFail(c)
		return
	}

	err := a.Srv.UpdateWorkOrder(ctx, contextx.FromUserID(ctx), workid, title, content)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithMessage(c, "修改成功")
}

func (a *WorkOrderApi) Create(c *gin.Context) {
	ctx := c.Request.Context()
	content := c.Request.PostFormValue("content")
	title := c.Request.PostFormValue("title")

	if len(content) == 0 || len(title) == 0 {
		ginx.ResFail(c)
		return
	}
	uid := contextx.FromUserID(ctx)
	fmt.Printf("uid: %v\n", uid)
	err := a.Srv.CreateWorkOrder(ctx, contextx.FromUserID(ctx), title, content)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithMessage(c, "添加成功")
}
func (a *WorkOrderApi) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Request.PostFormValue("userId")
	workOrderId := c.Request.PostFormValue("workOrderId")

	if len(id) == 0 || len(workOrderId) == 0 || contextx.FromUserID(ctx) != id {
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

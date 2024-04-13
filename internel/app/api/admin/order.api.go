package admin

import (
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

	res, err := a.Srv.LoadWorkListList(ctx, int(PageNo), int(PageSize), "", status)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, res)
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

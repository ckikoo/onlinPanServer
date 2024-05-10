package api

import (
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DingdanApi struct {
	Srv *service.DingdanService
}

func (api *DingdanApi) GetDingdanList(c *gin.Context) {
	ctx := c.Request.Context()

	no := c.PostForm("pageNo")
	size := c.PostForm("pageSize")
	if no == "" {
		no = "1"
	}
	if size == "" {
		size = "10"
	}

	no_, err := strconv.Atoi(no)
	if err != nil {
		ginx.ResFail(c)
		return
	}
	size_, err := strconv.Atoi(size)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	res, err := api.Srv.GetDingdanList(ctx, no_, size_, true, "")
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOkWithData(c, res)
}

func (api *DingdanApi) Buy(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.PostForm("id")
	if id == "" {
		ginx.ResFail(c)
	}

	_no, err := strconv.Atoi(id)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	err = api.Srv.Buy(ctx, contextx.FromUserID(ctx), _no)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOkWithMessage(c, "购买成功")
}

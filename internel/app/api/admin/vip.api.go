package admin

import (
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VipAPI struct {
	VipSrv *service.VipSrv
}

func (api *VipAPI) GetVipList(c *gin.Context) {

	pageNO := c.PostForm("pageNo")
	pageSize := c.PostForm("pageSize")

	if pageNO == "" {
		pageNO = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	_no, err := strconv.Atoi(pageNO)
	if err != nil {
		ginx.ResFailWithMessage(c, "参数错误")
		return
	}
	_size, err := strconv.Atoi(pageSize)

	if err != nil {
		ginx.ResFailWithMessage(c, "参数错误")
		return
	}

	res, err := api.VipSrv.GetVipList(c, _no, _size, "*")
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	ginx.ResOkWithData(c, res)
}

func (api *VipAPI) UpdateTime(c *gin.Context) {

	time := c.PostForm("time")
	id := c.PostForm("id")
	_time, err := strconv.Atoi(time)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	_id, err := strconv.Atoi(id)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	err = api.VipSrv.UpdateTime(c, _id, _time, contextx.FromUserID(c.Request.Context()))
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	ginx.ResOk(c)
}

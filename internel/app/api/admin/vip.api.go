package admin

import (
	"onlineCLoud/internel/app/dao/vip"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VipApi struct {
	Srv *service.VipService
}

func (api *VipApi) GetVipList(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := api.Srv.GetVipList(ctx)
	if err != nil {

		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOkWithData(c, res)
}
func (api *VipApi) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	show, err := strconv.ParseBool(c.PostForm("show"))
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	err = api.Srv.UpdateStatus(ctx, id, show)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOk(c)
}

func (api *VipApi) Update(c *gin.Context) {
	ctx := c.Request.Context()

	var body vip.Vip

	if err := ginx.ParseForm(c, &body); err != nil {
		ginx.ResFail(c)
		return
	}

	err := api.Srv.Update(ctx, body.ID, body)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOk(c)
}

func (api *VipApi) Add(c *gin.Context) {

	ctx := c.Request.Context()
	var vip vip.Vip
	if err := ginx.ParseForm(c, &vip); err != nil {
		ginx.ResFail(c)
		return
	}
	c1, _ := strconv.ParseBool(c.PostForm("show"))
	vip.Show = c1
	if err := api.Srv.AddVip(ctx, vip); err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOk(c)
}
func (api *VipApi) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		ginx.ResFail(c)
		return
	}

	err = api.Srv.DelVips(ctx, id)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOk(c)
}

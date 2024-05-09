package api

import (
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"

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

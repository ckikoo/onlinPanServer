package api

import (
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"

	"github.com/gin-gonic/gin"
)

type VipAPI struct {
	VipSrv *service.VipSrv
}

func (api *VipAPI) GetInfo(c *gin.Context) {

	res, err := api.VipSrv.GetInfo(c, contextx.FromUserID(c.Request.Context()))
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	ginx.ResOkWithData(c, res)
}

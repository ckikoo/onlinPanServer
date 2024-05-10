package api

import (
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"

	"github.com/gin-gonic/gin"
)

type PageApi struct {
	Srv *service.PageService
}

func (api *PageApi) GetPageList(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := api.Srv.GetPackageList(ctx, true)
	if err != nil {

		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOkWithData(c, res)
}

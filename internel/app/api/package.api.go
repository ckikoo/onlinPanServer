package api

import (
	"onlineCLoud/internel/app/dao/pkg"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/contextx"

	"github.com/gin-gonic/gin"
)

type PackageApi struct {
	Srv *service.PackageService
}

func (pack *PackageApi) GetPackInfo(c *gin.Context) {

	res, err := pack.Srv.GetPackInfo(c.Request.Context())
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOkWithData(c, res)

	return

}

func (pack *PackageApi) BuySpace(c *gin.Context) {
	ctx := c.Request.Context()

	sid := c.PostForm("packId")

	info, err := pack.Srv.CheckExists(ctx, sid)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ok, err := pack.Srv.BuySpace(ctx, contextx.FromUserID(ctx), info.(pkg.Pkg))
	if err != nil || !ok {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOk(c)
	return
}

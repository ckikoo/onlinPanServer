package admin

import (
	Package "onlineCLoud/internel/app/dao/package"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PackageApi struct {
	Srv *service.PageService
}

func (api *PackageApi) GetPackageList(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := api.Srv.GetPackageList(ctx, false)
	if err != nil {

		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOkWithData(c, res)
}
func (api *PackageApi) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.PostForm("ID"))
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

func (api *PackageApi) Update(c *gin.Context) {
	ctx := c.Request.Context()

	var body Package.Package

	if err := ginx.ParseForm(c, &body); err != nil {
		ginx.ResFail(c)
		return
	}

	err := api.Srv.Update(ctx, int(body.Model.ID), body)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	ginx.ResOk(c)
}

func (api *PackageApi) Add(c *gin.Context) {

	ctx := c.Request.Context()
	var Package Package.Package
	if err := ginx.ParseForm(c, &Package); err != nil {
		ginx.ResFail(c)
		return
	}
	c1, _ := strconv.ParseBool(c.PostForm("show"))
	Package.Show = c1
	if err := api.Srv.AddPackage(ctx, Package); err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOk(c)
}
func (api *PackageApi) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.PostForm("ID"))
	if err != nil {
		ginx.ResFail(c)
		return
	}

	err = api.Srv.DelPackages(ctx, id)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOk(c)
}

package admin

import (
	"onlineCLoud/internel/app/config"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminApi struct {
	AdminSrv *service.AdminSrv
}

func (a *AdminApi) LoadUserList(c *gin.Context) {
	ctx := c.Request.Context()
	pageNo := c.Request.PostFormValue("pageNo")
	pageSize := c.Request.PostFormValue("pageSize")
	nickNameFuzzy := c.PostForm("nickNameFuzzy")
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

	res, err := a.AdminSrv.LoadUserList(ctx, int(PageNo), int(PageSize), nickNameFuzzy, status)
	if err != nil {
		ginx.ResFail(c)
		return
	}

	ginx.ResOkWithData(c, res)
}

func (a *AdminApi) GetSysSettings(c *gin.Context) {

	res := map[string]interface{}{}
	res["userInitUseSpace"] = config.C.File.InitSpaceSize / 1024
	res["captchaLength"] = config.C.Captcha.Length
	res["downloadLimit"] = config.C.Download.Limit

	ginx.ResOkWithData(c, res)
}
func (a *AdminApi) SaveSysSettings(c *gin.Context) {

	captchaLength := c.PostForm("captchaLength")
	downloadLimitStr := c.PostForm("downloadLimit")
	userInitUseSpaceStr := c.PostForm("userInitUseSpace")

	downloadLimit, err := strconv.Atoi(downloadLimitStr)
	if err != nil {

		ginx.ResFailWithMessage(c, "nvalid download limit")
		return
	}

	// 将 userInitUseSpace 转换为布尔类型
	userInitUseSpace, err := strconv.Atoi(userInitUseSpaceStr)
	if err != nil {

		ginx.ResFailWithMessage(c, "nvalid userInitUseSpace value")
		return
	}
	capLength, err := strconv.Atoi(captchaLength)
	if err != nil {

		ginx.ResFailWithMessage(c, "nvalid userInitUseSpace value")
		return
	}

	if len(captchaLength) == 0 || downloadLimit == 0 || userInitUseSpace == 0 {
		// 错误处理：如果条件不符合预期，返回相应的错误信息
		ginx.ResFailWithMessage(c, "nvalid parameters")
		return
	}

	config.C.File.InitSpaceSize = uint64(userInitUseSpace)
	config.C.Captcha.Length = capLength
	config.C.Download.Limit = uint64(downloadLimit)
	config.SaveConfig("./config/config.bak.toml")
	ginx.ResOk(c)
}

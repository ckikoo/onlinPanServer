package router

import (
	"onlineCLoud/internel/app/api"
	"onlineCLoud/internel/app/middleware"
	"onlineCLoud/pkg/auth"

	"github.com/gin-gonic/gin"
)

type IRouter interface {
	Regitser(app *gin.Engine) error
	Prefixes() []string
}

type Router struct {
	Auth       auth.Auther
	LoginAPI   *api.LoginAPI
	UserApi    *api.UserAPI
	FileApi    *api.FileApi
	RecycleApi *api.RecycleApi
	ShareApi   *api.ShareApi
	AdminApi   *api.AdminApi
	PackageApi *api.PackageApi
}

func (a *Router) Regitser(app *gin.Engine) error {
	a.RegisterApI(app)
	return nil
}

func (a *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}

func (a *Router) RegisterApI(app *gin.Engine) {
	g := app.Group("/api")
	key := "online-cloud-server"
	g.Use(middleware.PrintUrlRequest())
	g.Use(middleware.SessionMW(key))
	g.Use(middleware.AuthMiddleware(a.Auth,
		middleware.AllowPathPrefixSkipper(
			"/api/login", "/api/checkCode",
			"/api/sendEmailCode", "/api/register",
			"/api/resetPwd", "/api/file/download/",
			"/api/showShare",
		)))

	g.Use(middleware.CORSMiddleware())

	g.GET("/checkCode", api.GenerateCaptcha)
	g.POST("/sendEmailCode", api.SendEmail)
	g.POST("/register", a.LoginAPI.Register)
	g.POST("/login", a.LoginAPI.Login)
	g.POST("/logout", a.LoginAPI.Logout)
	g.POST("/resetPwd", a.LoginAPI.ResetPasswd)
	g.POST("/updateUserAvatar", a.UserApi.UpdateUserAvatar)
	g.POST("/getUserInfo", a.UserApi.GetInfo)
	g.POST("/getUseSpace", a.UserApi.GetUserSpace)
	g.POST("/updatePassword", a.UserApi.UpdatePassword)
	g.GET("/getAvatar/:user", a.UserApi.GetUserAvatar)

	g.POST("/package/getinfo", a.PackageApi.GetPackInfo)
	g.POST("/package/buy", a.PackageApi.BuySpace)

	g.POST("/file/loadDataList", a.FileApi.GetFileList)
	g.POST("/file/uploadFile", a.FileApi.UploadFile)
	g.POST("/file/cancelUploadFile", a.FileApi.CancelUpload)
	g.POST("/file/newFoloder", a.FileApi.NewFoloder)
	g.GET("/file/getImage/:src", a.FileApi.GetImage)
	g.POST("/file/delFile", a.FileApi.DelFiles)
	g.GET("/file/ts/getVideoInfo/:fid", a.FileApi.GetVideoInfo)
	g.POST("/file/getFile/:fid", a.FileApi.GetFileInfo)
	g.GET("/file/getFile/:fid", a.FileApi.GetFileInfo)
	g.POST("/file/getFolderInfo", a.FileApi.GetFolderInfo)
	g.POST("/file/rename", a.FileApi.FileRename)
	g.POST("/file/changeFileFolder", a.FileApi.ChangeFileFolder)
	g.POST("/file/loadAllFolder", a.FileApi.LoadAllFolder)
	g.POST("/file/createDownloadUrl/:fid", a.FileApi.CreateDownloadUrl)
	g.GET("/file/download/:code", a.FileApi.Download)

	g.POST("/recycle/loadRecycleList", a.RecycleApi.GetFileList)
	g.POST("/recycle/recoverFile", a.RecycleApi.RecoverFile)
	g.POST("/recycle/delFile", a.RecycleApi.DelFiles)

	g.POST("/share/loadShareList", a.ShareApi.LoadShareList)
	g.POST("/share/shareFile", a.ShareApi.ShareFile)
	g.POST("/share/cancelShare", a.ShareApi.CancelShare)

	g.POST("/showShare/getShareLoginInfo", a.ShareApi.GetShareLoginInfo)
	g.POST("/showShare/getShareInfo", a.ShareApi.GetShareInfo)
	g.POST("/showShare/loadFileList", a.ShareApi.LoadFileList)
	g.POST("/showShare/getFolderInfo", a.ShareApi.GetFolderInfo)
	g.POST("/showShare/checkShareCode", a.ShareApi.CheckShareCode)
	//g.POST("/showShare/getFile/:id", a.ShareApi.GetFile)
	// TODO 需要完善接口
	//g.POST("/share/checkShareCode", a.ShareApi.CheckShareCode)

	g.POST("/admin/loadUserList", a.AdminApi.LoadUserList)
	g.POST("/admin/loadFileList", a.AdminApi.LoadFileList)
	g.POST("/admin/getFolderInfo", a.AdminApi.GetFolderInfo)
}

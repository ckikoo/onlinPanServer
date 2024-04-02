package app

import (
	"onlineCLoud/internel/app/config"
	"onlineCLoud/internel/app/middleware"
	"onlineCLoud/internel/app/router"

	"github.com/LyricTian/gzip"
	"github.com/gin-gonic/gin"
)

func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)
	app := gin.New()
	app.NoMethod(middleware.NoMethodHander())
	app.NoRoute(middleware.NoRouteHandler())

	// CORS
	if config.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	if config.C.GZIP.Enable {
		app.Use(gzip.Gzip(gzip.BestCompression,
			gzip.WithExcludedExtensions(config.C.GZIP.ExcludedExtentions),
			gzip.WithExcludedPaths(config.C.GZIP.ExcludedPaths),
		))
	}

	r.Regitser(app)

	return app
}

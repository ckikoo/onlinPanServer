package app

import (
	"onlineCLoud/pkg/auth"

	"github.com/gin-gonic/gin"
)

type Injector struct {
	Engine *gin.Engine
	Auth   auth.Auther
}

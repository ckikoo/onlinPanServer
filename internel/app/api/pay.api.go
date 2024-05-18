package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func GetImage(c *gin.Context) {
	f, _ := os.OpenFile("img/pay.jpg", os.O_RDONLY, os.ModePerm)
	http.ServeContent(c.Writer, c.Request, "pay.jpg", time.Now(), f)
}

package ginx

import (
	"fmt"
	"net/http"
	"onlineCLoud/pkg/errors"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	prefix     = "NetCloud"
	ReqBodyKey = prefix + "/req-body"
	ResBodyKey = prefix + "/res-body"
)

func GetToken(c *gin.Context) string {
	var token string
	auth, err := c.Cookie("jwt_token")
	if err != nil {
		return ""
	}
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

func GetBodyData(c *gin.Context) []byte {
	if v, ok := c.Get(ReqBodyKey); ok {
		if b, ok := v.([]byte); ok {
			return b
		}
	}
	return nil
}

func ParseParamID(c *gin.Context, key string) uint64 {
	val := c.Param(key)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// Parse body json data to struct
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("Parse request json failed: %s", err.Error()))
	}
	return nil
}

// Parse query parameter to struct
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.New(fmt.Sprintf("Parse request query failed: %s", err.Error()))
	}

	return nil
}

// Parse body form data to struct
func ParseForm(c *gin.Context, obj interface{}) error {

	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("Parse request form failed: %s", err.Error()))
	}
	fmt.Println("parse success")
	return nil
}

type Response struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	Info   string      `json:"info"`
	Status string      `json:"status"`
}

func ResData(c *gin.Context, code int, data []byte) {
	c.Data(code, "text/plain", data)
}

func ResJson(c *gin.Context, code int, data interface{}, msg string, status string) {
	c.JSON(http.StatusOK, Response{
		Code:   code,
		Data:   data,
		Info:   msg,
		Status: status,
	})
}
func ResNeedReload(c *gin.Context) {
	ResJson(c, 600, nil, "失效", "")
}
func ResOk(c *gin.Context) {
	ResJson(c, 200, nil, "操作成功", "success")
}

func ResOkWithMessage(c *gin.Context, message string) {
	ResJson(c, 200, nil, message, "success")
}

func ResOkWithData(c *gin.Context, data interface{}) {
	ResJson(c, 200, data, "查询成功", "success")
}

func ResFail(c *gin.Context) {
	ResJson(c, -1, nil, "操作失败", "fail")
}

func ResFailWithMessage(c *gin.Context, msg string) {
	ResJson(c, -1, nil, msg, "fail")
}

func ResFailWithData(c *gin.Context, data interface{}) {
	ResJson(c, -1, data, "操作失败", "fail")
}

func ResFailDetailed(c *gin.Context, data interface{}, msg string) {
	ResJson(c, -1, data, msg, "fail")
}

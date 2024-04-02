package service

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	captchaPre = "captcha_"
)

func Captcha(c *gin.Context, Type string, l, w, h int) {

	captchaId := captcha.NewLen(l)
	session := sessions.Default(c)
	key := fmt.Sprintf("%v%v", captchaPre, Type)
	session.Set(key, captchaId)
	_ = session.Save()
	_ = GenerateCaptchaHandler(c.Writer, c.Request, captchaId, w, h)
}

func GenerateCaptchaHandler(w http.ResponseWriter, r *http.Request, id string, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Paragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer

	_ = captcha.WriteImage(&content, id, width, height)
	http.ServeContent(w, r, id+".png", time.Time{}, bytes.NewReader(content.Bytes()))

	return nil
}

func CaptchaVerify(c *gin.Context, Type string, code string) bool {
	session := sessions.Default(c)
	key := fmt.Sprintf("%v%v", captchaPre, Type)

	if captchaId := session.Get(key); captchaId != nil {
		session.Delete(key)
		_ = session.Save()
		if captcha.VerifyString(captchaId.(string), code) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

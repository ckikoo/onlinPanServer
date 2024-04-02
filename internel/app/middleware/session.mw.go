package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func sessionConfig() sessions.Store {
	sessionMaxAge := 60 * 5
	sessionSecret := "session"
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		MaxAge: sessionMaxAge, // 秒
		Path:   "/",
	})
	return store
}

func SessionMW(keyPairs string) gin.HandlerFunc {
	store := sessionConfig()
	return sessions.Sessions(keyPairs, store)
}

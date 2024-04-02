package app

import (
	"onlineCLoud/internel/app/config"
	"onlineCLoud/pkg/auth"
	"onlineCLoud/pkg/auth/jwtauth"
	"onlineCLoud/pkg/auth/jwtauth/storer/buntdb"
	"onlineCLoud/pkg/auth/jwtauth/storer/redis"

	"github.com/dgrijalva/jwt-go"
)

func InitAuth() (auth.Auther, func(), error) {
	cfg := config.C.JWTAuth

	var opts []jwtauth.Option
	opts = append(opts, jwtauth.SetExpired(cfg.Expired))
	opts = append(opts, jwtauth.SetSigningKey([]byte(cfg.SigningKey)))
	opts = append(opts, jwtauth.SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(cfg.SigningKey), nil
	}))

	var method jwt.SigningMethod

	switch cfg.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}

	opts = append(opts, jwtauth.SetSigningMethod(method))

	var store jwtauth.Storer

	switch cfg.Store {
	case "redis":
		rcfg := config.C.Redis
		store = redis.NewStore(&redis.Config{
			Addr:      rcfg.Addr,
			Password:  rcfg.Password,
			DB:        cfg.RedisDB,
			Keyprefix: cfg.RedisPrefix,
		})
	default:

		s, err := buntdb.NewStore(cfg.FilePath)
		if err != nil {
			return nil, nil, err
		}
		store = s
	}

	auth := jwtauth.New(store, opts...)
	cleanFunc := func() {
		auth.Release()
	}

	return auth, cleanFunc, nil
}

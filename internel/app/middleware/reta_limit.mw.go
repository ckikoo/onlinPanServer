package middleware

import (
	"fmt"
	"onlineCLoud/internel/app/config"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/pkg/contextx"
	"onlineCLoud/pkg/errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"golang.org/x/time/rate"
)

func RateLimiterMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config.C.RateLimiter

	if !cfg.Enable {
		return EmptyMiddleware()
	}

	rc := config.C.Redis
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": rc.Addr,
		},
		Password: rc.Password,
		DB:       cfg.RedisDB,
	})

	limiter := redis_rate.NewLimiter(ring)
	limiter.Fallback = rate.NewLimiter(rate.Inf, 0)

	return func(ctx *gin.Context) {
		if SkipHandler(ctx, skippers...) {
			ctx.Next()
			return
		}

		userID := contextx.FromUserEmail(ctx.Request.Context())
		if userID != "" {
			limit := cfg.Count
			rate, delay, allow := limiter.AllowMinute(fmt.Sprintf("%v", userID), limit)
			if !allow {
				h := ctx.Writer.Header()
				h.Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
				h.Set("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
				delaySec := int64(delay / time.Second)
				h.Set("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))
				ginx.ResFailWithMessage(ctx, errors.ErrTooManyRequests)

				return
			}
		}
		ctx.Next()
	}
}

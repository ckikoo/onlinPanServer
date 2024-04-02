package cache

import (
	"context"
	"fmt"
	"onlineCLoud/internel/app/dao/redisx"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestCache(t *testing.T) {
	rdx := redis.NewClient(&redis.Options{
		Password: "123456",
		Addr:     "127.0.0.1:6379",
		DB:       0,
	})
	rd := redisx.NewClientWithClient(context.Background(), rdx)
	s := make([]string, 1)
	s[0] = "123456"
	rd.ZsetWithTimestamps(rdx.Context(), "debug", s, 1, time.Minute*30)
	rd.HSet(context.Background(), "debug1", s[0], "")

	m, err := rd.HGet(context.Background(), "debug1", "15")
	fmt.Printf("m: %v\n", len(m))

	fmt.Printf("err: %v\n", err == redis.Nil)

}

package test

import (
	"testing"

	"github.com/go-redis/redis"
)

func TestRedis(t *testing.T) {

	rd := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	// rd.Set("111", "111", time.Hour*time.Duration(10))
	cmd := rd.Get("user:space:lj_5683@163.com")
	if cmd.Err() != nil {
		panic(cmd.Err())
	}

	t.Log(cmd.Val())
}

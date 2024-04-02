package service

import (
	"context"
	"fmt"
	"onlineCLoud/internel/app/dao/mailx"
	"onlineCLoud/internel/app/dao/redisx"
	"time"

	"github.com/go-redis/redis"
)

func Email(ctx context.Context, dest string, code string) {
	rdx := redisx.NewClient()
	defer rdx.Close()
	go mailx.Email.SendMsgwithHtml(ctx, dest, "验证码", "验证码："+code)
	rdx.Set(ctx, "email@"+dest, code, time.Minute*15)
}

func CheckEmail(ctx context.Context, eamilAccount string, code string) (bool, error) {
	rdx := redisx.NewClient()
	defer rdx.Close()

	str, err := rdx.Get(ctx, "email@"+eamilAccount)
	if err != nil && err != redis.Nil {
		return false, err
	}

	return str == code, nil
}

func DeleteEmail(ctx context.Context, email string) {
	rdx := redisx.NewClient()
	defer rdx.Close()

	_, err := rdx.Delete(ctx, "email@"+email)
	if err != nil {
		fmt.Println("error when del email", email, err)

	}
}

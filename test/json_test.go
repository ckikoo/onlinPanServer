package test

import (
	"context"
	"fmt"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/user"
	"testing"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type aa struct {
	A int
	B int
}

// func TestJsonToMap(t *testing.T) {
// 	a := aa{1, 2}
// 	b := json.MarshalToString(a)
// 	fmt.Printf("b: %v\n", b)
// 	var c map[string]interface{}
// 	json.Unmarshal([]byte(b), &c)
// 	fmt.Printf("c: %v\n", c)
// }

func TestMapFormSql(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3308)/NetCloud?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "tb_", SingularTable: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	rd := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	ls := redisx.NewClientWithClient(context.Background(), rd)
	f := user.UserRepo{DB: db, Rd: ls}

	m := f.GetUseSpace(context.Background(), "lj_5683@163.com")
	fmt.Printf("m: %v\n", m)
}

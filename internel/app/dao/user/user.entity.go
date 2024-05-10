package user

import (
	"context"
	"onlineCLoud/internel/app/dao/util"
	"time"

	"gorm.io/gorm"
)

func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(User))
}

type User struct {
	UserID   string `json:"user_id" form:"user_id" gorm:"column:user_id;type:varchar(36) ;primaryKey"`
	Status   int8   `json:"status" form:"status" gorm:"column:status;type:tinyint(1);index:key_status"`
	NickName string `json:"nickName" form:"nickName" gorm:"column:nick_name;type:varchar(20);index:key_nick_name"`
	Password string `json:"password" form:"password" gorm:"column:password;type:varchar(32)"`
	Email    string `json:"email" form:"email" gorm:"column:email;type:varchar(30);uniqueIndex:key_email"`
	Avatar   string `json:"avatar" form:"avatar" gorm:"column:avatar;type:varchar(100)"`
	Admin    bool   `json:"admin" form:"admin" gorm:"column:admin"`
	VipId    int    `json:"vipId" form:"vipId" gorm:"column:vipId"`
	UserSpace
	CreateTime   time.Time `json:"joinTime"`
	LastJoinTime time.Time `json:"lastLoginTime"`
	EncPassWord  string    `form:"encPassword" gorm:"column:encPassword"`
}

type UserSpace struct {
	UseSpace   uint64 `json:"useSpace" form:"useSpace" gorm:"column:use_space;type:bigint(20) unsigned"`
	TotalSpace uint64 `json:"totalSpace" form:"totalSpace" gorm:"column:total_space;type:bigint(20) unsigned"`
}

func (u *UserSpace) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"useSpace":   u.UseSpace,
		"totalSpace": u.TotalSpace,
	}
}

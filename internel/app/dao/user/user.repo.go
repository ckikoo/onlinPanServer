package user

import (
	"context"
	"encoding/json"
	"log"
	"onlineCLoud/internel/app/dao/file"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/util"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/pkg/errors"
	logger "onlineCLoud/pkg/log"
	jsonutil "onlineCLoud/pkg/util/json"
	"time"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
	Rd *redisx.Redisx
}

var (
	User_KEY = "user_lock"
)

func (a *UserRepo) LoadUserList(ctx context.Context, p *schema.PageParams, userName string, status string) ([]User, error) {

	var list []User
	db := GetUserDB(ctx, a.DB)
	if userName != "" {
		db.Where("nick_name like ?", "%"+userName+"%")
	}
	if status != "*" {
		db.Where("status = ?", status)
	}

	err := util.WrapPageQuery(ctx, db, p, &list, true)

	return list, err
}

func (a *UserRepo) GetUserListTotal(ctx context.Context, p *schema.PageParams, userName string, status string) (int64, error) {

	db := GetUserDB(ctx, a.DB)
	if userName != "" {
		db.Where("nick_name like ?", "%"+userName+"%")
	}
	if status != "*" {
		db.Where("status = ?", status)
	}
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return total, err

}
func (a *UserRepo) FindOneByName(ctx context.Context, email string, out *User) error {

	db := GetUserDB(ctx, a.DB).Where("email = ?", email)
	ok, err := util.FindOne(ctx, db, out)

	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	return nil
}
func (a *UserRepo) FindOneById(ctx context.Context, id string, out *User) error {

	db := GetUserDB(ctx, a.DB).Where(&User{UserID: id})
	ok, err := util.FindOne(ctx, db, out)

	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	return nil
}
func (a *UserRepo) Create(ctx context.Context, item *User) error {
	result := GetUserDB(ctx, a.DB).Create(item)

	if result.Error != nil {
		return nil
	}

	return nil
}

func (a *UserRepo) Update(ctx context.Context, id string, item User) error {

	result := GetUserDB(ctx, a.DB).Where("user_id=?", id).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *UserRepo) UpdateUserStatus(ctx context.Context, id string, status int) error {

	result := GetUserDB(ctx, a.DB).Where("user_id=?", id).Update("status", status)
	if result.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return nil
}

func (a *UserRepo) Delete(ctx context.Context, id string) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Delete(User{})
	return errors.WithStack(result.Error)
}

func (a *UserRepo) UpdatePassword(ctx context.Context, email string, password string) error {
	return GetUserDB(ctx, a.DB).Where("email=?", email).Update("password", password).Error
}

func (a *UserRepo) FindAvatarByName(ctx context.Context, user_id string, out *string) error {
	db := GetUserDB(ctx, a.DB).Select("avatar").Where("user_id = ?", user_id)
	_, err := util.FindOne(ctx, db, out)
	if err != nil {
		return err
	}
	return nil
}
func (a *UserRepo) UpdateUserAvatar(ctx context.Context, email string, filename string) error {
	return GetUserDB(ctx, a.DB).Where("email = ?", email).UpdateColumn("avatar", filename).Error
}

func (a *UserRepo) SetRedis(ctx context.Context, email string, in string) error {
	return a.Rd.Set(ctx, email, in, time.Hour*(24))

}

func (a *UserRepo) GetUseSpace(ctx context.Context, email string) map[string]interface{} {

	space, _ := a.Rd.Get(ctx, "user:space:"+email)
	var item UserSpace
	if space != "" {
		json.Unmarshal([]byte(space), &item)
		return item.ToMap()
	}
	var useSpace int64

	res := file.GetFileDB(ctx, a.DB).Select("file_size").Joins("JOIN tb_user ON tb_user.user_id = tb_file.user_id AND tb_user.email = ?", email).Count(&useSpace)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		logger.Log("error", res.Error)
		return item.ToMap()
	}

	// 加上email
	err := GetUserDB(ctx, a.DB).Select("").Where(&User{Email: email})
	if err != nil {
		return item.ToMap()
	}

	str := jsonutil.MarshalToString(item)
	a.Rd.Set(ctx, "user:space:"+email, str, time.Hour*(24))
	return item.ToMap()
}

func (a *UserRepo) GetUserSpaceById(ctx context.Context, id string) UserSpace {

	var item User

	err := GetUserDB(ctx, a.DB).Select("use_space ", "total_space").Where(&User{UserID: id}).First(&item).Error
	if err != nil {
		return UserSpace{}
	}
	return item.UserSpace
}

func (a *UserRepo) UpdateSpace(ctx context.Context, email string, add uint64, update ...bool) error {
	a.Rd.Delete(ctx, "user:space:"+email)
	db := GetUserDB(ctx, a.DB).Where("email = ?", email)
	if len(update) == 0 {
		err := db.UpdateColumn("use_space", gorm.Expr("use_space + ?", add)).Error
		if err != nil {
			return err
		}
	} else {
		err := db.UpdateColumn("use_space", add).Error
		if err != nil {
			return err
		}
	}
	a.Rd.Delete(ctx, "user:space:"+email)
	return nil
}

func (a *UserRepo) UpdateEncPassword(ctx context.Context, email string, password string) error {

	db := GetUserDB(ctx, a.DB).Where("email = ?", email)

	err := db.UpdateColumn("encPassword", password).Error
	if err != nil {
		return err
	}
	return nil
}

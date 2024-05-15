package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/util"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/pkg/errors"
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
	space, err := a.Rd.Get(ctx, "user:space:"+email)
	var item UserSpace
	if err == nil {
		if err := json.Unmarshal([]byte(space), &item); err == nil {
			return item.ToMap()
		}
		log.Printf("Failed to unmarshal: %v", err)
	}

	var useSpace int64
	res := a.DB.Table("tb_file").Select("SUM(file_size) as space_used").
		Joins("JOIN tb_user ON tb_user.user_id = tb_file.user_id").
		Where("tb_user.email = ?", email).Scan(&useSpace)
	if res.Error != nil {
		log.Printf("Error querying use space: %v", res.Error)
		return item.ToMap()
	}
	fmt.Printf("useSpace: %v\n", useSpace)

	// Query for total space with COALESCE
	var totalSpace int64
	res = a.DB.Table("tb_user u").
		Select("COALESCE(p.spaceSize, u.total_space) AS space_size").
		Joins("JOIN tb_vip v ON u.user_id = v.user_id").
		Joins("JOIN tb_package p ON v.vip_package_id = p.id").
		Where("? BETWEEN v.active_from AND v.active_until AND u.email = ?", time.Now(), email).
		Order("COALESCE(p.spaceSize, u.total_space) DESC").
		Limit(1).
		Scan(&totalSpace)
	if res.Error != nil || res.RowsAffected == 0 {
		log.Printf("Error or no data for total space: %v", res.Error)
		return item.ToMap()
	}

	// Set the values and cache them
	item.TotalSpace = uint64(totalSpace)
	item.UseSpace = uint64(useSpace)
	str, _ := json.Marshal(item)
	a.Rd.Set(ctx, "user:space:"+email, str, 24*time.Hour)

	return item.ToMap()
}

func (a *UserRepo) GetUserSpaceById(ctx context.Context, id string) UserSpace {
	var item UserSpace
	var useSpace int64
	res := a.DB.Table("tb_file").Select("SUM(file_size) as space_used").
		Joins("JOIN tb_user ON tb_user.user_id = tb_file.user_id").
		Where("tb_file.user_id = ?", id).Scan(&useSpace)
	if res.Error != nil {
		log.Printf("Error querying use space: %v", res.Error)
		return item
	}

	var totalSpace int64
	res = a.DB.Table("tb_user u").
		Select("COALESCE(p.spaceSize, u.total_space) AS space_size").
		Joins("JOIN tb_vip v ON u.user_id = v.user_id").
		Joins("JOIN tb_package p ON v.vip_package_id = p.id").
		Where("? BETWEEN v.active_from AND v.active_until AND u.user_id = ?", time.Now(), id).
		Order("COALESCE(p.spaceSize, u.total_space) DESC").
		Limit(1).
		Scan(&totalSpace)
	if res.Error != nil || res.RowsAffected == 0 {
		log.Printf("Error or no data for total space: %v", res.Error)
		return item
	}

	item.TotalSpace = uint64(totalSpace)
	item.UseSpace = uint64(useSpace)

	return item
}

func (a *UserRepo) UpdateSpace(ctx context.Context, uid string, add uint64, update ...bool) error {
	a.Rd.Delete(ctx, "user:space:"+uid)
	db := GetUserDB(ctx, a.DB).Where(&User{UserID: uid})
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
	a.Rd.Delete(ctx, "user:space:"+uid)
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

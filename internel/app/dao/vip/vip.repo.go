package vip

import (
	"context"
	"errors"
	"onlineCLoud/internel/app/config"
	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/internel/app/dao/util"
	logger "onlineCLoud/pkg/log"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type VipRepo struct {
	VipDB *gorm.DB
}
type info struct {
	UserId      string    `gorm:"column:user_id" json:"user_id" form:"user_id"` // 页面名称
	VipId       int       `gorm:"column:vip_id" json:"vip_id" form:"vip_id"`    // 会员编号
	Email       string    `gorm:"column:email" json:"email" form:"email"`       // 邮箱
	NickName    string    `json:"nickName" form:"nickName" gorm:"column:nick_name;type:varchar(20);index:key_nick_name"`
	Avatar      string    `json:"avatar" form:"avatar" gorm:"column:avatar;type:varchar(100)"`
	PageName    string    `gorm:"column:pageName" json:"pageName" form:"pageName"` // 页面名称
	ActiveFrom  time.Time `json:"activeFrom"`
	ActiveUntil time.Time `json:"activeUntil"`
}

func (a *VipRepo) LoadVipInfoList(ctx context.Context, pageno, pageSize int, uid string) ([]info, error) {
	var list []info

	db := GetVipDb(ctx, a.VipDB).
		Select("tb_vip.id as vip_id,tb_user.nick_name, tb_user.avatar, tb_user.email , tb_package.pageName, tb_vip.user_id, tb_vip.active_from, tb_vip.active_until").
		Joins("JOIN tb_user ON tb_vip.user_id = tb_user.user_id").
		Joins("JOIN tb_package ON tb_vip.vip_package_id = tb_package.id").
		Limit(pageSize).Offset((pageno - 1) * pageSize).
		Find(&list)
	if db.Error != nil {
		return nil, db.Error
	}
	return list, nil
}

func (a *VipRepo) GetVipListTotal(ctx context.Context, pageno, pageSize int, uid string) (int64, error) {

	var total int64
	db := GetVipDb(ctx, a.VipDB).
		Select("tb_user.nick_name").
		Joins("JOIN tb_user ON tb_vip.user_id = tb_user.user_id").
		Joins("JOIN tb_package ON tb_vip.vip_package_id = tb_package.id").
		Limit(pageSize).Offset((pageno - 1) * pageSize).
		Count(&total)
	if db.Error != nil {

		return 0, nil
	}
	return total, nil

}

func (a *VipRepo) FindOneById(ctx context.Context, id uint, out *Vip) error {

	db := GetVipDb(ctx, a.VipDB).Where(&Vip{VipPackageID: id})
	ok, err := util.FindOne(ctx, db, out)

	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	return nil
}
func (a *VipRepo) Insert(ctx context.Context, in *Vip) error {
	_, err := redisx.NewClient().Delete(context.Background(), "vipInfo"+in.UserID)
	if err != nil && err != redis.Nil {
		logger.Log("Fatal", err)
	}

	return GetVipDb(ctx, a.VipDB).Create(in).Error
}

func (a *VipRepo) UpgradeExpireTime(userID string, days uint32) error {
	return a.VipDB.Transaction(func(tx *gorm.DB) error {
		currentMemberships := make([]Vip, 0)
		result := tx.Where("user_id = ? AND active_until > ?", userID, time.Now()).Order("active_until desc").Find(&currentMemberships)
		if result.Error != nil {
			return result.Error // 直接返回错误
		}

		for _, ship := range currentMemberships {
			ship.ActiveUntil = ship.ActiveUntil.Add(time.Hour * 24 * time.Duration(days))
			if err := tx.Save(&ship).Error; err != nil {
				return err // 在遇到错误时返回错误并回滚
			}
		}

		return nil // 确保所有操作成功后返回nil
	})
}
func (a *VipRepo) CheckExists(userID string, packid uint) (bool, error) {
	var item Vip
	res := GetVipDb(context.Background(), a.VipDB).Where(&Vip{UserID: userID, VipPackageID: packid}).First(&item)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, res.Error
	}

	// 如果查询成功且存在记录，则返回true
	if res.RowsAffected > 0 {
		return true, nil
	}

	// 默认返回false
	return false, nil
}
func (a *VipRepo) GetVipInfoByUserIDAndPackId(userID string, packid uint) (*Vip, error) {
	var item Vip
	res := GetVipDb(context.Background(), a.VipDB).Where(&Vip{UserID: userID, VipPackageID: packid}).First(&item)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}

	// 如果查询成功且存在记录，则返回true
	if res.RowsAffected > 0 {
		return &item, nil
	}

	// 默认返回false
	return nil, nil
}

func (a *VipRepo) GetVipInfoByUserID(userID string) (*vipInfo, error) {

	db := GetVipDb(context.Background(), a.VipDB)

	var info vipInfo

	res := db.Joins("JOIN tb_package ON tb_package.ID = tb_vip.vip_package_id").
		Where("tb_vip.user_id = ?", userID).
		Select("tb_vip.user_id, tb_vip.active_until, tb_package.pageName, tb_package.speedLimit, tb_package.spaceSize").
		Order("tb_package.spaceSize desc, tb_package.speedLimit desc").First(&info)

	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			return nil, res.Error
		}
		info = vipInfo{
			ActiveUntil: time.Now().Add(time.Hour * 24 * 9999),
			PageName:    "普通用户",
			SpeedLimit:  100,
			SpaceSize:   int(config.C.File.InitSpaceSize),
		}

	}

	return &info, nil
}

type vipInfo struct {
	UserID      string    `gorm:"column:user_id"`
	ActiveUntil time.Time `gorm:"column:active_until"`
	PageName    string    `gorm:"column:pageName"`
	SpeedLimit  int       `gorm:"column:speedLimit"`
	SpaceSize   int       `gorm:"column:spaceSize"`
}

func (a *VipRepo) UpdateTime(ctx context.Context, id, endTime int, uid string) error {
	db := GetVipDb(ctx, a.VipDB)
	e := time.Unix(int64(endTime), 0)
	res := db.Where(&Vip{UserID: uid, Model: gorm.Model{ID: uint(id)}}).Updates(&Vip{ActiveUntil: e})
	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}

		return errors.New("记录不存在")
	}
	return nil
}
func (a *VipRepo) Delete(ctx context.Context, id int) error {
	db := GetVipDb(ctx, a.VipDB)

	res := db.Where(&Vip{Model: gorm.Model{ID: uint(id)}}).Delete(nil)
	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}

		return errors.New("记录不存在")
	}
	return nil
}

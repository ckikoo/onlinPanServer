package vip

import (
	"context"
	"errors"
	"log"
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

func (a *VipRepo) LoadVipInfoList(ctx context.Context, uid string) ([]Vip, error) {

	list := make([]Vip, 0)
	res := GetVipDb(ctx, a.VipDB).Where(&Vip{UserID: uid}).Find(&list)
	if res.Error != nil {
		return nil, res.Error
	}
	return list, nil
}

func (a *VipRepo) GetVipListTotal(ctx context.Context, uid string) (int64, error) {

	db := GetVipDb(ctx, a.VipDB)

	var total int64
	err := db.Where(&Vip{UserID: uid}).Count(&total).Error
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return total, err

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

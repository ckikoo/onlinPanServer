package vip

import (
	"context"
	"errors"
	"log"
	"onlineCLoud/internel/app/schema"

	"gorm.io/gorm"
)

type VipRepo struct {
	DB *gorm.DB
}

func (f *VipRepo) GetVipList(ctx context.Context, page bool) ([]Vip, error) {
	db := GETVipDB(ctx, f.DB)

	temp := make([]Vip, 0)

	parms := schema.PageParams{
		PageNo:   0,
		PageSize: -1,
	}

	res := db.Offset(parms.PageNo).Limit(parms.PageSize).Find(&temp)
	if res.RowsAffected == 0 {
		return nil, nil
	}

	if res.Error != nil {
		return nil, res.Error
	}
	return temp, nil

}

func (f *VipRepo) GetVipListTotal(ctx context.Context) (int64, error) {
	db := GETVipDB(ctx, f.DB)

	var total int64

	err := db.Count(&total).Error
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return total, err
}

func (f *VipRepo) Insert(ctx context.Context, vip Vip) error {
	db := GETVipDB(ctx, f.DB)

	return db.Create(&vip).Error
}

// Delete removes multiple Vip records from the database by IDs.
func (f *VipRepo) DeleteBulk(ctx context.Context, ids []int) error {
	db := GETVipDB(ctx, f.DB)

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, id := range ids {
		if err := tx.Unscoped().Delete(&Vip{ID: id}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil { // 注意这里使用 tx.Commit()
		tx.Rollback()
		return err
	}

	return nil
}

// Delete removes multiple Vip records from the database by IDs.
func (f *VipRepo) Delete(ctx context.Context, id int) error {
	db := GETVipDB(ctx, f.DB)

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

		}
	}()

	if err := tx.Unscoped().Delete(&Vip{ID: id}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (f *VipRepo) Update(ctx context.Context, id int, vip Vip) error {
	db := GETVipDB(ctx, f.DB)
	res := db.Where("id = ?", id).Updates(&vip)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}
func (f *VipRepo) UpdateStatus(ctx context.Context, id int, status bool) error {
	db := GETVipDB(ctx, f.DB)
	res := db.Where("id = ?", id).Update("show", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

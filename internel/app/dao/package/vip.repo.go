package Pack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"onlineCLoud/internel/app/schema"

	"gorm.io/gorm"
)

type PackageRepo struct {
	DB *gorm.DB
}

func (f *PackageRepo) GetPageList(ctx context.Context, status bool, page bool) ([]Package, error) {
	db := GETPackageDB(ctx, f.DB)
	if status {
		fmt.Printf("status: %v\n", status)
		db = db.Where(&Package{Show: status})

	}
	temp := make([]Package, 0)

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

func (f *PackageRepo) GetPageListTotal(ctx context.Context, status bool) (int64, error) {
	db := GETPackageDB(ctx, f.DB)
	if status {
		fmt.Printf("status: %v\n", status)
		db = db.Where(&Package{Show: status})

	}
	var total int64

	err := db.Count(&total).Error
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return total, err
}

func (f *PackageRepo) Insert(ctx context.Context, Package Package) error {
	db := GETPackageDB(ctx, f.DB)

	return db.Create(&Package).Error
}

func (f *PackageRepo) Delete(ctx context.Context, id int) error {
	db := GETPackageDB(ctx, f.DB)

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

		}
	}()

	if err := tx.Unscoped().Where(&Package{Model: gorm.Model{ID: uint(id)}}).Delete(nil).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (f *PackageRepo) Update(ctx context.Context, id int, Package Package) error {
	db := GETPackageDB(ctx, f.DB)
	res := db.Where("id = ?", id).Updates(&Package)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}
func (f *PackageRepo) UpdateStatus(ctx context.Context, id int, status bool) error {
	db := GETPackageDB(ctx, f.DB)
	res := db.Where("id = ?", id).Update("show", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (f *PackageRepo) FindById(ctx context.Context, id int) (*Package, error) {
	db := GETPackageDB(ctx, f.DB)
	var temp Package
	res := db.Where("id = ?", id).First(&temp)
	if res.Error != nil && gorm.ErrRecordNotFound != res.Error {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	return &temp, nil
}

package dingdan

import (
	"context"
	"errors"
	"onlineCLoud/internel/app/schema"

	"gorm.io/gorm"
)

type DingdanRepo struct {
	DB *gorm.DB
}

func (f *DingdanRepo) GetDingdanList(ctx context.Context, pageNo, pageSize int, page bool, uid string) ([]Dingdan, error) {
	db := GetDingdanDB(ctx, f.DB)
	temp := make([]Dingdan, 0)

	params := schema.PageParams{
		PageNo:   pageNo,
		PageSize: pageSize,
	}

	res := db.Offset((params.PageNo - 1) * params.PageSize).Limit(params.PageSize).Order("created_at desc").Find(&temp)
	if res.Error != nil {
		return nil, res.Error
	}
	return temp, nil
}

func (f *DingdanRepo) GetDingdanListTotal(ctx context.Context, uid string) (int64, error) {
	db := GetDingdanDB(ctx, f.DB)

	db = db.Where(&Dingdan{UserId: uid})
	var total int64

	err := db.Model(&Dingdan{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (f *DingdanRepo) Insert(ctx context.Context, dingdan Dingdan) error {
	db := GetDingdanDB(ctx, f.DB)

	return db.Create(&dingdan).Error
}

func (f *DingdanRepo) Delete(ctx context.Context, id int) error {
	db := GetDingdanDB(ctx, f.DB)

	return db.Unscoped().Delete(&Dingdan{}, id).Error
}

func (f *DingdanRepo) Update(ctx context.Context, id int, dingdan Dingdan) error {
	db := GetDingdanDB(ctx, f.DB)

	res := db.Model(&Dingdan{}).Where("id = ?", id).Updates(&dingdan)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (f *DingdanRepo) FindById(ctx context.Context, id int) (*Dingdan, error) {
	db := GetDingdanDB(ctx, f.DB)

	var dingdan Dingdan
	res := db.First(&dingdan, id)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	return &dingdan, nil
}

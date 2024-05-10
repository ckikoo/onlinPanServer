package dingdan

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type DingdanRepo struct {
	DB *gorm.DB
}

type dingdan struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	PackName  string  `gorm:"column:packageName" json:"packageName"`
	NickName  string  `gorm:"column:nick_name" json:"nick_name"`
	UserId    string  `gorm:"column:user_id" json:"user_id"`
	PackageId int     `gorm:"column:package_id" json:"package_id"`
	Price     float32 `gorm:"column:price" json:"price" form:"price"`
}

func (f *DingdanRepo) GetDingdanList(ctx context.Context, pageNo, pageSize int, page bool, uid string) ([]dingdan, error) {
	db := GetDingdanDB(ctx, f.DB)
	temp := make([]dingdan, 0)

	res := db.
		Select("tb_dingdan.id, tb_dingdan.created_at, tb_package.pageName as packageName, tb_dingdan.user_id as user_id, tb_dingdan.package_id as PackageId, tb_package.price, tb_user.nick_name").
		Joins("left join tb_package on tb_dingdan.package_id = tb_package.id").
		Joins("left join tb_user on tb_dingdan.user_id = tb_user.user_id").
		Where("tb_user.nick_name LIKE ?", "%"+uid+"%").
		Offset((pageNo - 1) * pageSize).
		Limit(pageSize).
		Order("tb_dingdan.created_at desc").
		Scan(&temp)
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

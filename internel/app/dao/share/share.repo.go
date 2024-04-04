package share

import (
	"context"
	"log"
	"onlineCLoud/internel/app/dao/util"
	"onlineCLoud/internel/app/schema"
	"onlineCLoud/pkg/errors"
	"strings"

	"gorm.io/gorm"
)

type ShareRepo struct {
	DB *gorm.DB
}

func (f *ShareRepo) GetShareList(ctx context.Context, uid string, schema *schema.RequestFileListPage, page bool) ([]Share, error) {
	db := GetShareDB(ctx, f.DB)

	if schema.OrderBy != "" {
		db.Order(schema.OrderBy)
	}
	sql := "SELECT tb_file.file_name, tb_share.* FROM tb_share LEFT JOIN tb_file ON tb_file.file_id = tb_share.file_id WHERE tb_share.user_id = ?"

	db = db.Raw(sql, uid)

	var temp []Share
	err := util.WrapPageQuery(ctx, db, &schema.PageParams, &temp, page)

	return temp, err

}

func (f *ShareRepo) GetShareListTotal(ctx context.Context, uid string, schema *schema.RequestFileListPage) (int64, error) {
	db := GetShareDB(ctx, f.DB)

	db = db.Where("user_id=?", uid)
	var total int64

	err := db.Count(&total).Error
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return total, err
}

func (f *ShareRepo) Insert(ctx context.Context, share Share) error {
	db := GetShareDB(ctx, f.DB)

	return db.Create(&share).Error
}

func (f *ShareRepo) CancelShare(ctx context.Context, uid string, fids []string) error {
	db := GetShareDB(ctx, f.DB)
	db.Begin()
	// 使用原始 SQL 表达式执行删除操作
	sql := "DELETE FROM tb_share WHERE user_id = ? AND FIND_IN_SET(share_id, ?)"
	res := db.Exec(sql, uid, strings.Join(fids, ","))
	if res.Error != nil || int(res.RowsAffected) != len(fids) {
		db.Rollback()
		return res.Error
	}

	return nil
}
func (f *ShareRepo) GetShareInfo(ctx context.Context, shareId string) (*Share, error) {
	db := GetShareDB(ctx, f.DB)
	share := new(Share)
	res := db.Where(&Share{ShareId: shareId}).First(share)
	if res.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	return share, nil
}
func (f *ShareRepo) UpdateShareShowCount(ctx context.Context, shareId string) (bool, error) {
	db := GetShareDB(ctx, f.DB)
	res := db.Where("id = ?", shareId).UpdateColumn("show_count", gorm.Expr("show_count + ?", 1))
	if res.RowsAffected == 0 {
		return false, errors.New("not found")
	}
	return true, nil
}

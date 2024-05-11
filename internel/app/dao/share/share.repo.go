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
type share struct {
	ShareId    string `json:"shareId" form:"shareId" gorm:"column:share_id;type:varchar(36) ;primaryKey"`
	FileId     string `json:"fileId" form:"fileId" gorm:"column:file_id;type:varchar(36) ;"`
	UserId     string `json:"userId" form:"userId" gorm:"column:user_id;type:varchar(36) ;"`
	ExpireTime string `json:"expireTime"`
	ShareTime  string `json:"shareTime"`
	Code       string `json:"code" form:"code" gorm:"column:code;type:varchar(5) ;"`
	ValidType  int8   `json:"validType" form:"validType" gorm:"column:valid_type;type:tinyint(1);"`
	ShowCount  uint32 `json:"showCount" form:"showCount" gorm:"column:show_count;type:int(11);"`
	FileName   string `json:"file_name"   form:"file_name" gorm:"column:file_name"` //文件名
}

func (f *ShareRepo) GetShareList(ctx context.Context, uid string, schema *schema.RequestFileListPage, page bool) ([]share, error) {
	db := GetShareDB(ctx, f.DB)

	if schema.OrderBy != "" {
		db.Order(schema.OrderBy)
	}
	sql := "SELECT tb_file.*, tb_share.* FROM tb_share LEFT JOIN tb_file ON tb_file.file_id = tb_share.file_id WHERE tb_share.user_id = ?"

	db = db.Raw(sql, uid)

	var temp []share
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
	// 使用原始 SQL 表达式执行删除操作
	sql := "DELETE FROM tb_share WHERE user_id = ? AND FIND_IN_SET(share_id, ?)"
	res := db.Exec(sql, uid, strings.Join(fids, ","))
	if res.Error != nil || int(res.RowsAffected) != len(fids) {
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
func (f *ShareRepo) AddShareShowCount(ctx context.Context, shareId string) (bool, error) {
	db := GetShareDB(ctx, f.DB)
	res := db.Where("id = ?", shareId).UpdateColumn("show_count", gorm.Expr("show_count + ?", 1))
	if res.RowsAffected == 0 {
		return false, errors.New("not found")
	}
	return true, nil
}

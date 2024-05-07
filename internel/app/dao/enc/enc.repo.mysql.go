package enc

import (
	"context"
	"errors"
	"log"
	"onlineCLoud/internel/app/dao/util"
	"onlineCLoud/internel/app/schema"

	"gorm.io/gorm"
)

type EncRepo struct {
	Db *gorm.DB
}

func (repo *EncRepo) AddFile(ctx context.Context, enc *Enc) error {
	return repo.Db.WithContext(ctx).Create(enc).Error
}

func (f *EncRepo) GetFileList(ctx context.Context, uid string, schema *schema.RequestFileListPage, page bool) ([]Enc, error) {
	db := GetEncDB(ctx, f.Db)

	if schema.FilePid != "" {
		db = db.Where("file_pid=?", schema.FilePid)
	}
	if schema.FolderType != 0 {
		db = db.Where("folder_type=?", schema.FolderType)
	}
	if schema.Path != nil && len(schema.Path) > 0 {
		db = db.Where("file_id in (?)", schema.Path)
	}

	// 排除 //防止自己覆盖自己
	if schema.ExInclude != nil && len(schema.ExInclude) > 0 {
		db = db.Where("file_id not in (?)", schema.ExInclude)
	}
	if schema.OrderBy != "" {
		db.Order(schema.OrderBy)
	}
	db = db.Where("user_id=?", uid)

	db.Order("folder_type desc")
	db.Order("create_time asc")
	var temp []Enc
	err := util.WrapPageQuery(ctx, db, &schema.PageParams, &temp, page)

	return temp, err
}

func (f *EncRepo) GetFileListTotal(ctx context.Context, uid string, schema *schema.RequestFileListPage) (int64, error) {
	db := GetEncDB(ctx, f.Db)

	if schema.FilePid != "" {
		db = db.Where("file_pid=?", schema.FilePid)
	}
	if schema.FolderType != 0 {
		db = db.Where("folder_type=?", schema.FolderType)
	}
	if schema.Path != nil && len(schema.Path) > 0 {
		db = db.Where("file_id in (?)", schema.Path)
	}

	// 排除 //防止自己覆盖自己
	if schema.ExInclude != nil && len(schema.ExInclude) > 0 {
		db = db.Where("file_id not in (?)", schema.ExInclude)
	}
	if schema.OrderBy != "" {
		db.Order(schema.OrderBy)
	}
	db = db.Where("user_id=?", uid)
	var total int64

	err := db.Count(&total).Error
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return total, err
}

func (f *EncRepo) AddEncFile(ctx context.Context, file *Enc) error {
	err := GetEncDB(ctx, f.Db).Create(file).Error
	if err != nil {
		return err
	}

	return nil

}
func (f *EncRepo) InsertFileBatch(ctx context.Context, file []Enc) error {
	err := GetEncDB(ctx, f.Db).CreateInBatches(file, len(file)).Error
	if err != nil {
		return err
	}

	return nil

}
func (f *EncRepo) DelFiles(ctx context.Context, uid string, fileId []string) error {
	if len(fileId) == 0 {
		return nil
	}
	db := GetEncDB(ctx, f.Db)

	return db.Where("user_id=?", uid).Delete("file_id in (?)", fileId).Error
}
func (f *EncRepo) GetFileInfo(ctx context.Context, fileId string, uid string) (*Enc, error) {
	db := GetEncDB(ctx, f.Db)

	var file Enc
	if err := db.Where("user_id=?", uid).Where("file_id =?", fileId).First(&file).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &file, nil
}

func (f *EncRepo) CheckFileName(ctx context.Context, filePId string, uid string, fileName string) (*Enc, error) {
	db := GetEncDB(ctx, f.Db)

	var file Enc
	db = db.Where("user_id=?", uid).Where("file_name=?", fileName)
	if filePId != "" {
		db = db.Where("file_pid=?", filePId)
	}

	if err := db.First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &file, nil
		} else {
			return nil, err

		}
	}
	return &file, nil
}

func (f *EncRepo) FileRename(ctx context.Context, uid string, fileId string, filePid string, fileName string) (bool, error) {
	db := GetEncDB(ctx, f.Db)

	db = db.Where("user_id=?", uid).Where("file_id=?", fileId).Where("file_pid=?", filePid).Update("file_name", fileName)
	if db.Error != nil {
		return false, db.Error
	}
	if db.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (f *EncRepo) UpdateFile(ctx context.Context, file *Enc) error {
	db := GetEncDB(ctx, f.Db)

	return db.Where("user_id=?", file.UserID).Where("file_id=?", file.FileID).UpdateColumns(file).Error
}

package file

import (
	"context"
	"errors"
	"fmt"
	"log"
	"onlineCLoud/internel/app/dao/util"
	"onlineCLoud/internel/app/schema"

	"gorm.io/gorm"
)

type FileRepo struct {
	Db *gorm.DB
}

func (f *FileRepo) GetFileList(ctx context.Context, uid string, schema *schema.RequestFileListPage, page bool) ([]File, error) {
	db := GetFileDB(ctx, f.Db)

	if schema.Category != "" && schema.Category != "all" {
		db = db.Where("file_category =?", schema.Category)
	}

	if schema.FileNameFuzzy != "" {
		db = db.Where("file_name like ?", "%"+schema.FileNameFuzzy+"%")
	}

	if schema.FilePid != "" {
		db = db.Where("file_pid=?", schema.FilePid)
	}
	if schema.FolderType != 0 {
		db = db.Where("folder_type=?", schema.FolderType)
	}
	if schema.Path != nil && len(schema.Path) > 0 {
		db = db.Where("file_id in (?)", schema.Path)
	}

	db = db.Where("del_flag = ?", schema.DelFlag)

	// 排除 //防止自己覆盖自己
	if schema.ExInclude != nil && len(schema.ExInclude) > 0 {
		db = db.Where("file_id not in (?)", schema.ExInclude)
	}
	if schema.OrderBy != "" {
		db.Order(schema.OrderBy)
	}
	if uid != "*" {
		db = db.Where("user_id=?", uid)
	}

	db.Order("folder_type desc")
	db.Order("create_time asc")
	var temp []File
	err := util.WrapPageQuery(ctx, db, &schema.PageParams, &temp, page)

	return temp, err
}

func (f *FileRepo) GetTotalUseSpace(ctx context.Context, uid string, out *uint64) error {
	db := GetFileDB(ctx, f.Db)
	return db.Where("user_id = ? ", uid).Pluck("COALESCE(SUM(file_size), 0) as total", out).Error
}

func (f *FileRepo) GetFileListTotal(ctx context.Context, uid string, schema *schema.RequestFileListPage) (int64, error) {
	db := GetFileDB(ctx, f.Db)

	if schema.Category != "" && schema.Category != "all" {
		db = db.Where("file_category =?", schema.Category)
	}

	if schema.FileNameFuzzy != "" {
		db = db.Where("file_name like ?", "%"+schema.FileNameFuzzy+"%")
	}

	if schema.FilePid != "" {
		db = db.Where("file_pid=?", schema.FilePid)
	}
	if schema.FolderType != 0 {
		db = db.Where("folder_type=?", schema.FolderType)
	}
	if schema.Path != nil && len(schema.Path) > 0 {
		db = db.Where("file_id in (?)", schema.Path)
	}
	db = db.Where("del_flag = ?", schema.DelFlag)
	db = db.Where("user_id=?", uid)
	var total int64

	err := db.Count(&total).Error
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return total, err
}

func (f *FileRepo) CheckFileExists(ctx context.Context, md5 string) (*File, error) {

	file, err := FileCache.FindMd5(ctx, md5)
	if err == CachaNoFound {

		db := GetFileDB(ctx, f.Db)
		db = db.Where("file_md5 =?", md5)

		var file File
		err := db.First(&file).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		FileCache.AddFile(ctx, file)
		return &file, nil
	} else {
		return file, err
	}

}

func (f *FileRepo) UploadFile(ctx context.Context, file *File) error {
	err := GetFileDB(ctx, f.Db).Create(file).Error
	if err != nil {
		return err
	}

	FileCache.AddFile(ctx, *file)
	return nil

}
func (f *FileRepo) DelFiles(ctx context.Context, uid string, fileId []string) error {
	if len(fileId) == 0 {
		return nil
	}
	db := GetFileDB(ctx, f.Db)

	fmt.Printf("fileId: %v\n", fileId)

	return db.Where("user_id=?", uid).Delete("file_id in (?)", fileId).Error

}
func (f *FileRepo) GetFileInfo(ctx context.Context, fileId string, uid string) (*File, error) {
	db := GetFileDB(ctx, f.Db)

	var file File
	if err := db.Where("user_id=?", uid).Where("file_id =?", fileId).First(&file).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &file, nil
}

func (f *FileRepo) CheckFileName(ctx context.Context, filePId string, uid string, fileName string, folderType string) (*File, error) {
	db := GetFileDB(ctx, f.Db)

	var file File
	db = db.Where("user_id=?", uid).Where("file_name=?", fileName)
	if filePId != "" {
		db = db.Where("file_pid=?", filePId)
	}
	if folderType != "" {
		db = db.Where("folder_type=?", folderType)
	}

	if err := db.First(&file).Error; err != nil {
		fmt.Printf("err: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("???")
			return &file, nil
		} else {
			return nil, err

		}
	}
	return &file, nil
}

func (f *FileRepo) FileRename(ctx context.Context, uid string, fileId string, filePid string, fileName string) (bool, error) {
	db := GetFileDB(ctx, f.Db)

	db = db.Where("user_id=?", uid).Where("file_id=?", fileId).Where("file_pid=?", filePid).Update("file_name", fileName)
	if db.Error != nil {
		return false, db.Error
	}
	if db.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (f *FileRepo) UpdateFile(ctx context.Context, file *File) error {
	db := GetFileDB(ctx, f.Db)

	return db.Where("user_id=?", file.UserID).Where("file_id=?", file.FileID).UpdateColumns(file).Error
}

func (f *FileRepo) UpdateFileDelFlag(ctx context.Context, UserID string, filePids []string, fileIds []string, oldFlag, newFlag int8, reTime string) error {
	db := GetFileDB(ctx, f.Db)
	if filePids != nil && len(filePids) > 0 {
		db = db.Where("file_pid in (?)", filePids)
	}
	if fileIds != nil && len(fileIds) > 0 {
		db = db.Where("file_id in (?)", fileIds)
	}
	return db.Where("user_id=?", UserID).Where("del_flag=?", oldFlag).Updates(map[string]interface{}{"del_flag": newFlag, "recovery_time": reTime}).Error
}

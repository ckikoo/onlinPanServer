package file

import (
	"context"
	"errors"
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

	db = db.Where("secure=?", schema.Secure)
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
	db = db.Where("secure=?", schema.Secure)
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

func (f *FileRepo) UploadFile(ctx context.Context, file *File) error {
	err := GetFileDB(ctx, f.Db).Create(file).Error
	if err != nil {
		return err
	}

	return nil

}
func (f *FileRepo) InsertFileBatch(ctx context.Context, file []File) error {
	err := GetFileDB(ctx, f.Db).CreateInBatches(file, len(file)).Error
	if err != nil {
		return err
	}

	return nil

}
func (f *FileRepo) DelFiles(ctx context.Context, uid string, fileId []string) error {
	if len(fileId) == 0 {
		return nil
	}
	db := GetFileDB(ctx, f.Db)

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

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &file, nil
		} else {
			return nil, err

		}
	}
	return &file, nil
}

func (f *FileRepo) CountFileByMd5(ctx context.Context, md5 string) (int64, error) {
	var count int64

	db := GetFileDB(ctx, f.Db)

	db = db.Where(File{
		FileMd5: md5,
	})

	if err := db.Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			return 0, err

		}
	}

	return count, nil
}

func (f *FileRepo) FindFilesByMd5s(ctx context.Context, md5 []string) ([]File, error) {

	db := GetFileDB(ctx, f.Db)
	files := make([]File, 0)
	if err := db.Where("file_md5 IN (?)", md5).Group("file_md5").Find(&files).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return files, nil
}

func (f *FileRepo) GetFileByMd5(ctx context.Context, md5 string) (*File, error) {
	db := GetFileDB(ctx, f.Db)

	var file File
	db = db.Where(File{
		FileMd5: md5,
	})

	if err := db.First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
func (f *FileRepo) UpdateFileSecure(ctx context.Context, UserID string, fileIds []string, status bool) error {
	db := GetFileDB(ctx, f.Db)

	if fileIds != nil && len(fileIds) > 0 {
		db = db.Where("file_id in (?)", fileIds)
	}
	return db.Where("user_id=?", UserID).Updates(map[string]interface{}{"secure": status}).Error
}

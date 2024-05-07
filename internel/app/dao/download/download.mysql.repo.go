package download

import (
	"context"

	"gorm.io/gorm"
)

type DownloadRepo struct {
	Db *gorm.DB
}

func (repo *DownloadRepo) CreateRecord(ctx context.Context, data Download) (bool, error) {
	db := getDownloadB(ctx, repo.Db)

	db = db.Create(&data)
	if db.RowsAffected != 0 {
		return true, nil
	}

	return false, db.Error

}

func (repo *DownloadRepo) FindRecordByCode(ctx context.Context, code string) (string, error) {
	db := getDownloadB(ctx, repo.Db)

	var result struct {
		FilePath string `gorm:"column:path"`
	}

	if err := db.Table("tb_download").
		Select("tb_file.file_path as path").
		Joins("JOIN tb_file ON tb_file.fild_id = tb_download.fild_id and tb_file.user_id = tb_download.user_id").
		Where("tb_download.code = ?", code).
		First(&result).Error; err != nil {
		return "", err
	}

	return result.FilePath, nil
}

func (repo *DownloadRepo) GetAllData(ctx context.Context) ([]Download, error) {
	db := getDownloadB(ctx, repo.Db)
	download := make([]Download, 0)

	err := db.Find(&download).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return download, nil
}

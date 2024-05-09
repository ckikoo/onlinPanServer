package download

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type DownloadRepo struct {
	Db *gorm.DB
}

func (repo *DownloadRepo) CreateRecord(ctx context.Context, data *Download) (bool, error) {
	db := getDownloadB(ctx, repo.Db)

	db = db.Create(data)
	fmt.Printf("data: %v\n", data)
	if db.RowsAffected != 0 {
		return true, nil
	}

	return false, db.Error

}

func (repo *DownloadRepo) FindRecordByCode(ctx context.Context, code string) (string, error) {
	db := getDownloadB(ctx, repo.Db)

	var resultMap map[string]interface{}
	if err := db.Table("tb_download").
		Select("tb_file.file_path as path").
		Joins("JOIN tb_file ON tb_file.file_id = tb_download.file_id AND tb_file.user_id = tb_download.user_id").
		Where("tb_download.code = ?", code).
		Scan(&resultMap).Error; err != nil {
		return "", err
	}

	filePath, ok := resultMap["path"].(string)
	if !ok {
		return "", fmt.Errorf("file_path not found in the result")
	}

	return filePath, nil
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

func (repo *DownloadRepo) Delete(ctx context.Context, code string) error {
	db := getDownloadB(ctx, repo.Db)

	res := db.Where(&Download{Code: code}).Delete(nil)
	return res.Error
}

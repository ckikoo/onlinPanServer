package download

import (
	"context"
	"onlineCLoud/internel/app/dao/util"
	"strconv"

	"gorm.io/gorm"
)

type Download struct {
	UserId     string `gorm:"column:user_id;index"`
	FileId     string `gorm:"column:file_id;index"`
	Code       string `gorm:"column:code;index"`
	CreateTime int64  `gorm:"column:create_time"`
}

func getDownloadB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Download))
}

func ToMap(data Download) map[string]interface{} {
	dataMap := map[string]interface{}{
		"code":        data.Code,
		"create_time": data.CreateTime,
		"user_id":     data.UserId,
		"file_id":     data.FileId,
	}
	return dataMap
}

func MapToStruct(dataMap map[string]string) Download {
	var download Download
	download.Code = dataMap["code"]
	timer, _ := strconv.ParseInt(dataMap["create_time"], 10, 64)
	download.CreateTime = timer
	download.UserId = dataMap["user_id"]
	download.FileId = dataMap["file_id"]
	return download
}

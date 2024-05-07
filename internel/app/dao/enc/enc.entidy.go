package enc

import (
	"context"
	"onlineCLoud/internel/app/dao/util"
	"time"

	"gorm.io/gorm"
)

type Enc struct {
	EncID      string    `json:"encId"     form:"encId" gorm:"column:encId;type:varchar(36);primaryKey"`                // 编号
	FileID     string    `json:"fileId"    form:"fileId" gorm:"column:file_id;type:varchar(36);primaryKey"`             // 文件编号
	UserID     string    `json:"userId"    form:"userId" gorm:"column:user_id;type:varchar(36);primaryKey"`             // 用户编号
	FolderType int8      `json:"folderType" form:"folderType" gorm:"column:folder_type;type:tinyint(1)"`                // 是否目录
	FileSize   uint64    `json:"fileSize"  form:"fileSize" gorm:"column:file_size;type:bigint(20);index:key_file_size"` // 文件大小
	FileMd5    string    `json:"fileMd5"   form:"fileMd5" gorm:"column:file_md5;type:varchar(32);index:key_file_md5"`   // 文件MD5妙传
	FilePid    string    `json:"filePid"   form:"filePid" gorm:"column:file_pid;type:varchar(36)"`                      // 文件父级pid
	FileName   string    `json:"fileName"  form:"fileName" gorm:"column:file_name;type:varchar(255)"`                   // 文件名
	FileCover  string    `json:"fileCover" form:"fileCover" gorm:"column:file_cover;type:varchar(255);"`                // 封面
	FilePath   string    `json:"filePath"  form:"filePath" gorm:"column:file_path;type:varchar(255)"`                   // 文件路径
	JoinTime   time.Time `json:"joinTime"  gorm:"column:join_time;type:datetime"`                                       // 加入时间
}

func GetEncDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Enc))
}

package file

import (
	"context"
	"onlineCLoud/internel/app/dao/util"
	"onlineCLoud/pkg/util/json"
	"time"

	"gorm.io/gorm"
)

func GetFileDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(File))
}

type File struct {
	FileID         string    `json:"fileId"     form:"fileId" gorm:"column:file_id;type:varchar(36) ;primaryKey;"`           //文件编号
	UserID         string    `json:"userId"     form:"userId" gorm:"column:user_id;type:varchar(36);primaryKey"`             //用户编号
	FolderType     int8      `json:"folderType" form:"folderType" gorm:"column:folder_type;type:tinyint(1)"`                 //是否目录
	FileType       int8      `json:"fileType"   form:"fileType" gorm:"column:file_type;type:tinyint(1)"`                     //文件类型（视频，音频，图片，pdf，doc，exec， 7 txt 8. code 9 zip 10 other ）
	FileCategory   int8      `json:"fileCategory" form:"fileCategory" gorm:"column:file_category;type:tinyint(1)"`           //视频音频图片文档其他  文件分类 1 开始
	Status         int8      `json:"status"      form:"status" gorm:"column:status;type:tinyint(1)" `                        //状态  转码中，失败，成功
	DelFlag        int8      `json:"delFlag"    form:"delFlag" gorm:"column:del_flag;size:tiny(1);index:key_del_flag"`       //
	FileSize       uint64    `json:"fileSize"   form:"fileSize" gorm:"column:file_size;type:bigint(20);index:key_file_size"` //文件大小
	FileMd5        string    `json:"fileMd5"    form:"fileMd5" gorm:"column:file_md5;type:varchar(32);index:key_file_md5"`   //文件MD5妙传
	FilePid        string    `json:"filePid"    form:"filePid" gorm:"column:file_pid;type:varchar(36)"`                      //文件父级pid
	FileName       string    `json:"fileName"   form:"fileName" gorm:"column:file_name;type:varchar(255)"`                   //文件名
	FileCover      string    `json:"fileCover"  form:"fileCover" gorm:"column:file_cover;type:varchar(255);"`                //封面
	FilePath       string    `json:"filePath"   form:"filePath" gorm:"column:file_path;type:varchar(255)"`                   // 文件路径
	CreateTime     time.Time `json:"createTime"   form:"createTime" gorm:"column:create_time"`                               //创建时间
	LastUpdateTime time.Time `json:"lastUpdateTime" form:"lastUpdateTime" gorm:"column:last_update_time"`                    //上一次访问时间
	Secure         bool
	JoinTime       time.Time `json:"join_time" form:"join_time" gorm:"column:join_time"` // 加入密码箱时间
}

func File2file() {

}

func Sfile2File() {

}

func FileToMap(file File) map[string]interface{} {
	fileMap := make(map[string]interface{})
	fileMap["file_id"] = file.FileID
	fileMap["user_id"] = file.UserID
	fileMap["folder_type"] = int(file.FolderType)
	fileMap["file_type"] = int(file.FileType)
	fileMap["file_category"] = int(file.FileCategory)
	fileMap["status"] = int(file.Status)
	fileMap["del_flag"] = int(file.DelFlag)
	fileMap["file_size"] = file.FileSize
	fileMap["file_md5"] = file.FileMd5
	fileMap["file_pid"] = file.FilePid
	fileMap["file_name"] = file.FileName
	fileMap["file_cover"] = file.FileCover
	fileMap["file_path"] = file.FilePath

	// 将时间转换为字符串格式
	fileMap["create_time"] = file.CreateTime
	fileMap["last_update_time"] = file.LastUpdateTime

	fileMap["secure"] = file.Secure
	fileMap["join_time"] = file.JoinTime

	return fileMap
}
func ToMd5Map(files []File) map[string]interface{} {
	md5Map := make(map[string]interface{})

	for _, file := range files {
		md5Map[file.FileMd5] = json.MarshalToString(file) // 将文件 MD5 作为键，文件结构体作为值
	}

	return md5Map
}

package recycle

import (
	"context"
	"onlineCLoud/internel/app/define"
	"onlineCLoud/internel/app/schema"
	"time"

	"gorm.io/gorm"
)

type RecycleRepo struct {
	DB *gorm.DB
}

type recycle struct {
	FileID         string `json:"fileId"     form:"fileId" gorm:"column:file_id;type:varchar(36) ;primaryKey;"`           //文件编号
	UserID         string `json:"userId"     form:"userId" gorm:"column:user_id;type:varchar(36);primaryKey"`             //用户编号
	FolderType     int8   `json:"folderType" form:"folderType" gorm:"column:folder_type;type:tinyint(1)"`                 //是否目录
	FileType       int8   `json:"fileType"   form:"fileType" gorm:"column:file_type;type:tinyint(1)"`                     //文件类型（视频，音频，图片，pdf，doc，exec， 7 txt 8. code 9 zip 10 other ）
	FileCategory   int8   `json:"fileCategory" form:"fileCategory" gorm:"column:file_category;type:tinyint(1)"`           //视频音频图片文档其他  文件分类 1 开始
	Status         int8   `json:"status"      form:"status" gorm:"column:status;type:tinyint(1)" `                        //状态  转码中，失败，成功
	DelFlag        int8   `json:"delFlag"    form:"delFlag" gorm:"column:del_flag;size:tiny(1);index:key_del_flag"`       //
	FileSize       uint64 `json:"fileSize"   form:"fileSize" gorm:"column:file_size;type:bigint(20);index:key_file_size"` //文件大小
	FileMd5        string `json:"fileMd5"    form:"fileMd5" gorm:"column:file_md5;type:varchar(32);index:key_file_md5"`   //文件MD5妙传
	FilePid        string `json:"filePid"    form:"filePid" gorm:"column:file_pid;type:varchar(36)"`                      //文件父级pid
	FileName       string `json:"fileName"   form:"fileName" gorm:"column:file_name;type:varchar(255)"`                   //文件名
	FileCover      string `json:"fileCover"  form:"fileCover" gorm:"column:file_cover;type:varchar(255);"`                //封面
	FilePath       string `json:"filePath"   form:"filePath" gorm:"column:file_path;type:varchar(255)"`                   // 文件路径
	CreateTime     string `json:"createTime"   form:"createTime" gorm:"column:create_time"`                               //创建时间
	LastUpdateTime string `json:"lastUpdateTime" form:"lastUpdateTime" gorm:"column:last_update_time"`                    //上一次访问时间
	Secure         bool
	JoinTime       string `json:"join_time" form:"join_time" gorm:"column:join_time"`           // 加入密码箱时间
	RecoveryTime   string `json:"recoveryTime" form:"recoveryTime" gorm:"column:recovery_time"` //进入回收站时间
}

func (repo *RecycleRepo) GetFileList(ctx context.Context, uid string, page schema.PageParams, isPage bool) ([]recycle, error) {
	db := GetRecycleDB(ctx, repo.DB)
	baseQuery := `
        SELECT 
            f.*, 
            r.recovery_time 
        FROM 
            File f 
        LEFT JOIN 
            Recycle r 
        ON 
            f.file_id = r.file_id AND f.user_id = r.user_id 
        WHERE 
            f.user_id = ? AND f.del_flag = ? 
    `
	var query string
	var args []interface{}

	if isPage {
		query = baseQuery + "LIMIT ? OFFSET ?;"
		args = append(args, uid, define.FileFlagInRecycleBin, page.GetPageSize(), (page.GetCurrentPage()-1)*page.GetPageSize())
	} else {
		query = baseQuery
		args = append(args, uid, define.FileFlagInRecycleBin)
	}

	var temp []recycle
	res := db.Raw(query, args...).Scan(&temp)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	return temp, nil
}

func (repo *RecycleRepo) GetFileListTotal(ctx context.Context, uid string) (int64, error) {
	db := GetRecycleDB(ctx, repo.DB)
	var total int64
	query := `
        SELECT 
            COUNT(*) 
        FROM 
            File f 
        LEFT JOIN 
            Recycle r 
        ON 
            f.file_id = r.file_id AND f.user_id = r.user_id 
        WHERE 
            f.user_id = ? AND f.del_flag = ?;
    `
	res := db.Raw(query, uid, define.FileFlagInRecycleBin).Scan(&total)
	if res.Error != nil {
		return 0, res.Error
	}
	return total, nil
}

func (repo *RecycleRepo) Delete(ctx context.Context, uid string, fileId string) error {
	db := GetRecycleDB(ctx, repo.DB)

	res := db.Where("user_id = ? AND file_id = ?", uid, fileId).Delete(&Recycle{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *RecycleRepo) Add(ctx context.Context, uid string, fileId string) error {
	db := GetRecycleDB(ctx, repo.DB)

	res := db.Create(&Recycle{UserID: uid, FileID: fileId, RecoveryTime: time.Now().Format("2006-01-02 15:04:05")})

	return res.Error
}

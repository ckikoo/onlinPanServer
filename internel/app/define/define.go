package define

import "strings"

const (
	FileShare1Day = iota
	FileShare7Day
	FileShare30Day
	FileShareForverDay
)

// （视频，音频，图片，pdf，doc，exec， 7 txt 8. code 9 zip 10 other ）
const (
	FileTypeFolder = iota
	FileTypeVideo
	FileTypeMusic
	FileTypeImage
	FileTypePDF
	FileTypeDoc
	FileTypeExcel
	FileTypeTxt
	FileTypeCode
	FileTypeZip
	FileTypeOther
)

// 文件仓库定义
const (
	FileCategoryFolder = iota
	FileCategoryVideo
	FileCategoryMusic
	FileCategoryImage
	FileCategoryDoc
	FileCategoryOthers
)

// 文件状态 (转码中，失败，成功)
const (
	FileStatusTraning = iota
	FileStatusFail
	FileStatusUsing
)

// 文件标记
const (
	FileFlagInUse        = iota //文件在使用
	FileFlagInRecycleBin        // 进入回收站
	FileFlagSoftDeleted         // 文件进入隐藏
)

const (
	DownloadStatusSuccess = iota
	DownloadStatusClientQuick
)

func GetDay(Dt int8) int {
	switch Dt {
	case FileShare1Day:
		return 1
	case FileShare7Day:
		return 7
	case FileShare30Day:
		return 30
	case FileShareForverDay:
		return -1
	default:
		return -2
	}
}

func GetFileType(filename string) int8 {

	switch filename {
	case "mp4", "avi", "mkv":
		return FileTypeVideo
	case "mp3", "wav", "flac":
		return FileTypeMusic
	case "jpg", "jpeg", "png", "gif":
		return FileTypeImage
	case "pdf":
		return FileTypePDF
	case "doc", "docx":
		return FileTypeDoc
	case "xls", "xlsx":
		return FileTypeExcel
	case "txt":
		return FileTypeTxt
	case "go", "java", "py", "cpp":
		return FileTypeCode
	case "zip", "rar":
		return FileTypeZip
	default:
		return FileTypeOther
	}
}

func FileCategoryID4Str(fileType int8) string {
	switch fileType {
	// 视频音频图片文档其他
	case FileCategoryFolder:
		return "folder"
	case FileCategoryVideo:
		return "video"
	case FileCategoryMusic:
		return "music"
	case FileCategoryImage:
		return "image"
	case FileCategoryDoc:
		return "doc"
	default:
		return "others"
	}
}
func FileCategoryStr4ID(fileType string) int8 {
	fileType = strings.ToLower(fileType)
	switch fileType {

	case "folder":
		return FileCategoryFolder
	case "video", "mp4", "avi", "mkv":
		return FileCategoryVideo
	case "music", "mp3", "wav", "flac":
		return FileCategoryMusic
	case "image", "jpg", "jpeg", "png", "gif":
		return FileCategoryImage
	case "doc", "docx", "txt", "md", "pdf":
		return FileCategoryDoc
	default:
		return FileCategoryOthers
	}
}

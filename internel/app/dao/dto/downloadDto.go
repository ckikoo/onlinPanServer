package dto

import "time"

type DownloadDto struct {
	Code     string
	FileName string
	Path     string
	FileSize uint64
	Modi     time.Time
}

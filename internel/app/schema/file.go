package schema

type FileUpload struct {
	FileId     string `form:"fileId"`
	FileName   string `form:"fileName"`
	FilePid    string `form:"filePid"`
	FileMd5    string `form:"fileMd5"`
	ChunkIndex int    `form:"chunkIndex"`
	Chunks     int    `form:"chunks"`
	FileSize   int    `form:"fileSize"`
}

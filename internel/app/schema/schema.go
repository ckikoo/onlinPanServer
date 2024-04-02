package schema

type PageParams struct {
	PageNo   int `form:"pageNo,default=0" json:"pageNo" `
	PageSize int `form:"pageSize,default=20" json:"pageSize"`
}

func (p PageParams) GetCurrentPage() int {
	return p.PageNo
}

func (p PageParams) GetPageSize() int {
	pageSize := p.PageSize
	if p.PageSize <= 0 {
		pageSize = 10
	}
	return pageSize
}

type RequestFileListPage struct {
	PageParams
	FileNameFuzzy string `form:"fileNameFuzzy" json:"fileNameFuzz"`
	Category      string `form:"category" json:"category"`
	FilePid       string `form:"filePid" json:"filePid"`
	Path          []string
	ExInclude     []string
	OrderBy       string
	FolderType    int8
	DelFlag       int8
}

type ListResult struct {
	Parms      *PageParams
	TotalCount int64       `json:"totalCount"`
	PageTotal  int64       `json:"pageTotal"`
	List       interface{} `json:"list"`
}

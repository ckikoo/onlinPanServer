package schema

type ShareInfo struct {
	ShareTime   string      `json:"shareTime"`
	ExpireTime  string      `json:"expireTime"`
	NickName    string      `json:"nickName"`
	FileName    string      `json:"fileName"`
	CurrentUser bool        `json:"currentUser"`
	FileID      string      `json:"fileId"`
	Avatar      interface{} `json:"avatar"`
	UserID      string      `json:"userId"`
}

type RequestShareListPage struct {
	PageParams
	ShareId string `json:"shareId" form:"shareId"`
	FilePid string `json:"filePid" form:"filePid"`
}

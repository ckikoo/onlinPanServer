package json

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// 定义JSON操作
var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

// MarshalToString JSON编码为字符串
func MarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return ""
	}
	return s
}

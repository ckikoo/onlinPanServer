package define

// 文件分享的日期
const (
	FileShare1Day = iota
	FileShare7Day
	FileShare30Day
	FileShareForverDay
)

func GetShareDay(Dt int8) int {
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

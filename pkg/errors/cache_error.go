package errors

var (
	CachaNoFound = New("cache not found")
	CacheIsEmpty = New("cache empty")
)

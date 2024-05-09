package download

import (
	"context"
	"fmt"
	"onlineCLoud/internel/app/dao/redisx"
	"strings"
	"time"
)

var (
	downloadKey = "onlinePanServer:download:%v"
)

func CreateRecordUseRedis(ctx context.Context, data Download, path string) (bool, error) {
	rd := redisx.NewClient()
	defer rd.Close() // Close the Redis client when the function exits

	diff := time.Unix(data.CreateTime, 0).Add(time.Hour * 1).Sub(time.Now())
	err := rd.Set(ctx, fmt.Sprintf(downloadKey, data.Code), path, diff)
	if err != nil {
		return false, err
	}

	return true, nil
}

func FindRecordByCode(ctx context.Context, code string) (string, error) {
	rd := redisx.NewClient()
	defer rd.Close() // Close the Redis client when the function exits

	path, err := rd.Get(ctx, fmt.Sprintf(downloadKey, code))
	if err != nil {
		return "", err
	}

	pos := strings.Index(path, ".")
	if pos == -1 {
		return "", fmt.Errorf("invalid path format")
	}

	return path, nil
}

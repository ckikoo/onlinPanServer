package download

import (
	"context"
	"fmt"
	"onlineCLoud/internel/app/dao/redisx"
	"time"
)

var (
	downloadKey = "onlinePanServer:download:%v"
)

func CreateRecordUseRedis(ctx context.Context, data Download, path string) (bool, error) {

	rd := redisx.NewClient()
	diff := time.Unix(data.CreateTime, 0).Add(time.Hour * 24).Sub(time.Now())
	err := rd.Set(ctx, fmt.Sprintf(downloadKey, data.Code), path, diff)
	if err != nil {
		return false, err
	}

	return true, nil

}

func FindRecordByCode(ctx context.Context, code string) (string, error) {
	rd := redisx.NewClient()

	path, err := rd.Get(ctx, fmt.Sprintf(downloadKey, code))
	if err != nil {
		return "", err
	}

	return path, nil
}
